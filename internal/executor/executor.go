package executor

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/models"
	"goosefs-cli2api/pkg/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/noop/log"

	"github.com/google/uuid"
)

type TaskRequest struct {
	TaskName string               `json:"task_name"`
	Command  string               `json:"command"`
	Action   models.GooseFSAction `json:"action" binding:"required"` // 入库用
	Path     string               `json:"path" binding:"required"`   // 入库用
	Args     []string             `json:"args"`
}

// 只允许内部调用，不允许外部传入所有指令，防止执行影响系统的指令
func addTask(req TaskRequest) (string, error) {
	taskID := uuid.New().String()
	cmd := exec.Command(req.Command, req.Args...)

	outputPath := utils.GenerateTaskID(req.TaskName, taskID)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	cmd.Stdout = outFile
	cmd.Stderr = outFile

	err = cmd.Start()
	if err != nil {
		return "", err
	}
	// 创建 task
	err = config.DB.CreateGoosefsTask(taskID, models.GoosefsTaskRequest{
		TaskName: &req.TaskName,
		Path:     &req.Path,
		Action:   req.Action,
	})
	if err != nil {
		return "", fmt.Errorf("createGoosefsTask error: %v", err)
	}

	go func() {
		cmd.Wait()
		outFile.Close()
		// 读取文件行数
		var count int
		defer func() {
			// 任务结束后更新任务状态
			err := config.DB.UpdateGoosefsTask(taskID, models.UpdateGoosefsTaskRequest{ExitCode: tea.String(cmd.ProcessState.String()), Count: tea.Int(count)})
			if err != nil {
				log.Errorf("UpdateGoosefsPathStatus error: %v", err)
			}
		}()
		file, err := os.Open(outputPath)
		if err != nil {
			count = -1
			log.Errorf("count error: open file error: %v", err)
			return
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			count = -2
			log.Errorf("count error: read file error: %v", err)
			return
		}
		lines := strings.Split(string(bytes), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			count++
		}
	}()

	return taskID, nil
}

func runCmd(cmd string, args []string) (string, error) {
	cmdObj := exec.Command(cmd, args...)
	bytes, err := cmdObj.Output()
	return string(bytes), err
}

func GetTaskStatus(filter models.FilterGoosefsTaskRequest) (models.TasksStatus, error) {
	tasks, err := config.DB.GetGoosefsTask(filter)
	if err != nil {
		return models.TasksStatus{}, err
	}
	if len(tasks) == 0 {
		return models.TasksStatus{}, fmt.Errorf("no task found by filter: %s", tea.Prettify(filter))
	}
	resp := models.TasksStatus{
		Data:   make(map[string]models.TaskInfo, len(tasks)),
		Status: models.TaskStatusSuccess,
	}
	successCount := 0
	isRunning := false
	for _, task := range tasks {
		// Path必须会有的，这里防止错误情况
		if task.Path == nil {
			continue
		}
		//缓存目录没有变化的不展示出来,当请求为GFSForceLoad,GFSDistributeLoad时
		if task.Count != nil && *task.Count == 0 && (task.Action == models.GFSForceLoad || task.Action == models.GFSDistributeLoad) {
			continue
		}

		taskinfo := models.TaskInfo{
			Path: *task.Path,
		}

		if task.ExitCode != nil && task.Count != nil {
			taskinfo.Count = *task.Count
			taskinfo.ExitCode = *task.ExitCode
			// 任务执行完成
			if GetCmdStatus(*task.ExitCode) == models.TaskStatusSuccess {
				successCount++
			}
		}
		if task.ExitCode == nil {
			isRunning = true
		}

		resp.Data[task.ID] = taskinfo
	}
	if isRunning {
		resp.Status = models.TaskStatusRunning
	} else {
		if successCount != len(tasks) {
			if successCount > 0 {
				resp.Status = models.TaskStatusNotallSuccess
			} else {
				resp.Status = models.TaskStatusFailed
			}
		}
	}

	return resp, nil
}

/*
status只返回成功Successfully内容的服务，其他 path的不通知
$ cat test_task_name_160d7bf1-0da3-473d-b4c4-29f7c7d37e14.txt
Allow up to 100 active jobs
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/20240813_075219_00032_byuem-d4444027-b4ef-4cf5-b3c8-6141b9884e93 loading
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/00000000000000000000.json loading
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/_trino_meta/extended_stats.json loading
Successfully loaded path /data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/20240813_075219_00032_byuem-d4444027-b4ef-4cf5-b3c8-6141b9884e93 after 1 attempts
Successfully loaded path /data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/00000000000000000000.json after 1 attempts
Successfully loaded path /data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/_trino_meta/extended_stats.json after 1 attempts

$ cat test_task_name_cd3fe749-4e32-4415-85df-57c3ddbea1b4.txt
Allow up to 100 active jobs
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/20240813_075219_00032_byuem-d4444027-b4ef-4cf5-b3c8-6141b9884e93 is already fully loaded in GooseFS
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/00000000000000000000.json is already fully loaded in GooseFS
/data-datalake-dataprod-bj-1251949819/deltalake/npd_temp.db/ods_corpdata_pingan_tb_certificate_integrate/_delta_log/_trino_meta/extended_stats.json is already fully loaded in GooseFS
*/
// 支持多个任务的结果一起查询输出
func GetTaskOutput(req models.QueryTaskRequest) (map[string]string, error) {
	// 获取 ID 反解析出原始的 TaskID
	taskFiles, err := utils.FindFiles(req)
	if err != nil {
		return nil, err
	}
	log.Debugf("taskFiles: %v", taskFiles)
	contentAll := make(map[string]string, len(taskFiles))
	for _, taskFile := range taskFiles {
		taskid, err := utils.ParseTaskID(taskFile)
		if err != nil {
			return nil, err
		}
		outputPath := strings.TrimSuffix(*config.Config.OutputDir, "/") + "/" + taskFile
		content, err := ioutil.ReadFile(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("taskID with taskName %s not found in output dir: %s", taskid, outputPath)
			}
			contentAll[taskid] = err.Error()
			continue
		}
		contentAll[taskid] = string(content)
	}

	return contentAll, nil
}
