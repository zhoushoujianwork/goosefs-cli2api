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
	"sync"

	"github.com/xops-infra/noop/log"

	"github.com/google/uuid"
)

type TaskStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type TaskRequest struct {
	Name    string   `json:"name"` // 因为会参与到文件名，所以这里不用最好用英文且不支持空格等
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

var (
	tasks      = make(map[string]*exec.Cmd)
	tasksMutex sync.RWMutex
)

func StartTaskManager() {
	// 可以扩展任务清理逻辑等
}

// 只允许内部调用，不允许外部传入所有指令，防止执行影响系统的指令
func addTask(req TaskRequest) (string, error) {
	taskID := uuid.New().String()
	cmd := exec.Command(req.Command, req.Args...)
	tasksMutex.Lock()
	tasks[taskID] = cmd
	tasksMutex.Unlock()

	outputPath := utils.GenerateTaskID(req.Name, taskID)
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

	go func() {
		cmd.Wait()
		outFile.Close()
	}()

	return taskID, nil
}

func GetTaskStatus(req models.QueryTaskRequest) (models.TaskStatus, error) {
	taskIDs := []string{}
	// 先判断依据有 2 个信息的情况，不实用遍历路径方式，降低复杂度
	if req.TaskID != nil && *req.TaskID != "" && req.TaskName == nil {
		taskIDs = append(taskIDs, *req.TaskID)
	} else {
		taskFiles, err := utils.FindFiles(req)
		if err != nil {
			return models.TaskStatus{}, fmt.Errorf("find files error: %v", err)
		}
		if len(taskFiles) == 0 {
			return models.TaskStatus{}, os.ErrNotExist
		}
		for _, taskFile := range taskFiles {
			taskid, err := utils.ParseTaskID(taskFile)
			if err != nil {
				return models.TaskStatus{}, err
			}
			taskIDs = append(taskIDs, taskid)
		}
	}

	resp := models.TaskStatus{
		Data:   make(map[string]string, len(taskIDs)),
		Status: models.TaskStatusSuccess,
	}
	for _, taskID := range taskIDs {
		tasksMutex.RLock()
		defer tasksMutex.RUnlock()
		cmd, exists := tasks[taskID]
		if !exists {
			resp.Data[taskID] = fmt.Sprintf("task %s not found in mem, maybe server has restart, you can call output api to get task output", taskID)
			resp.Status = models.TaskStatusFailed
		} else {
			resp.Data[taskID] = fmt.Sprintf("cmd: %s status: %s", strings.Join(cmd.Args, " "), cmd.ProcessState.String())
			if cmd.ProcessState.String() == "<nil>" {
				resp.Status = models.TaskStatusRunning
			} else if cmd.ProcessState.String() != "exit status 0" {
				resp.Status = models.TaskStatusFailed
			}
		}
	}
	return resp, nil

}

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
