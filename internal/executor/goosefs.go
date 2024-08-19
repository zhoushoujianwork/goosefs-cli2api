package executor

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/models"
	"goosefs-cli2api/pkg/dingtalk"
	"strings"
	"time"

	"github.com/xops-infra/noop/log"

	"github.com/alibabacloud-go/tea/tea"
)

// 实现任务完成后的告警通知，只检查大量任务带有 task_name 的，单独不带task_name的不做通知
func checkTasksIsFinished(act models.GooseFSAction, task_name *string) {
	if task_name == nil {
		log.Debugf("task_name and taskids are empty, skip checkTasksIsFinished")
		return
	}

	for {
		time.Sleep(5 * time.Second)
		log.Debugf("checkTasksIsFinished")

		status, err := GetTaskStatus(models.FilterGoosefsTaskRequest{
			TaskName: task_name,
			Action:   &act,
		})
		if err != nil {
			log.Errorf("checkTasksIsFinished error: %s", err)
			continue
		}

		var msg string
		if status.Status != models.TaskStatusSuccess {
			msg = "告警:" + string(act) + " " + string(status.Status) + " for task " + tea.StringValue(task_name) + "\n"
		} else {
			msg = "通知:" + string(act) + "task " + tea.StringValue(task_name) + "" + "\tSUCCESS!\n"
		}

		dingtalk.SendAlert(msg + tea.Prettify(status))

		break
	}
}

/*
./bin/goosefs fs distributedLoad --replication 1 /data-datalake/deltalake/aaa.db/bbb/

支持多路径 Path 任务提交，返回 1 个 task_id
实现任务完成告警，每次任务进来后挂起一个任务查询的协程，直到任务状态不是 running；
*/

func DistrubuteLoad(req models.GooseFSRequest) ([]string, error) {
	taskids := make([]string, 0)
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return nil, fmt.Errorf("path is required, should not be empty")
		}
		taskID, err := addTask(TaskRequest{
			Action:   models.GFSDistributeLoad,
			Path:     *p,
			TaskName: tea.StringValue(req.TaskName),
			Command:  *config.Config.Bin,
			Args:     []string{"fs", "distributedLoad", "--replication", "1", *p},
			// Args: []string{"fs", "distributedLoad", "--replication", "1", *p, "|grep 'Successfully loaded path'"}, // 只存入新加载的文件其他无关信息过滤掉
		})
		if err != nil {
			return nil, err
		}
		taskids = append(taskids, taskID)
	}
	go checkTasksIsFinished(models.GFSDistributeLoad, req.TaskName)
	return taskids, nil
}

/*
./bin/goosefs fs loadMetadata -R /data-datalake/deltalake/aaa.db/bbb/
*/
func LoadMetadata(req models.GooseFSRequest) ([]string, error) {
	taskids := make([]string, 0)
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return nil, fmt.Errorf("path is required, should not be empty")
		}
		taskID, err := addTask(TaskRequest{
			TaskName: tea.StringValue(req.TaskName),
			Command:  *config.Config.Bin,
			Args:     []string{"fs", "loadMetadata", "-R", *p},
			Path:     *p,
			Action:   models.GFSLoadMetadata,
		})
		if err != nil {
			return nil, err
		}
		taskids = append(taskids, taskID)
	}
	go checkTasksIsFinished(models.GFSLoadMetadata, req.TaskName)
	return taskids, nil
}

/*
GooseFSForceLoad 该步骤执行的是先去 LoadMetadata，然后再去 DistributeLoad，这样彻底更新
先执行 LoadMetadata，成功后再执行 DistributeLoad
*/
func ForceLoad(req models.GooseFSRequest) ([]string, error) {
	taskids := make([]string, 0)
	var errs []string
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return nil, fmt.Errorf("path is required, should not be empty")
		}
		_, err := runCmd(*config.Config.Bin, []string{"fs", "loadMetadata", "-R", *p})
		if err != nil {
			log.Errorf("loadMetadata error: %s", err)
			errs = append(errs, fmt.Errorf("loadMetadata for %s error: %s", *p, err).Error())
			continue
		}
		taskID, err := addTask(TaskRequest{
			TaskName: tea.StringValue(req.TaskName),
			Command:  *config.Config.Bin,
			Args:     []string{"fs", "distributedLoad", "--replication", "1", *p},
			Path:     *p,
			Action:   models.GFSForceLoad,
		})
		if err != nil {
			return nil, err
		}
		taskids = append(taskids, taskID)
	}
	if len(errs) > 0 {
		if len(taskids) != 0 {
			return nil, fmt.Errorf("success ids: %s, some task is failed: %s", strings.Join(taskids, ","), strings.Join(errs, "\n"))
		}
		return nil, fmt.Errorf(fmt.Sprintf("all task failed: %s", strings.Join(errs, "\n")))
	}
	go checkTasksIsFinished(models.GFSForceLoad, req.TaskName)
	return taskids, nil
}

/*
实时输出
./bin/goosefs fs ls /data-datalake/deltalake/aaa.db
*/
func List(path string, timeOut int) (string, error) {
	taskid, err := addTask(TaskRequest{
		Command: *config.Config.Bin,
		Args:    []string{"fs", "ls", path},
		Path:    path,
		Action:  models.GFSList,
	})
	if err != nil {
		return "", err
	}
	// log.Infof("taskid:", taskid)
	// wait for task done
	count := 0
	for {
		count++
		if count > timeOut {
			return "", fmt.Errorf("wait for task done timeout, you can call output api to get task output, taskid: %s", taskid)
		}
		// log.Infof("get task status:", taskid)
		status, err := GetTaskStatus(models.FilterGoosefsTaskRequest{
			TaskID: &taskid,
		})
		if err != nil {
			return "", err
		}
		log.Infof("task %s status: %s\n", taskid, status.Status)
		if status.Status == "<nil>" || status.Status == models.TaskStatusRunning {
			time.Sleep(1 * time.Second)
		} else if status.Status == models.TaskStatusSuccess {
			break
		} else {
			return "", fmt.Errorf("task %s exec error: %s", taskid, status.Status)
		}
	}

	// 读取输出文件
	output, err := GetTaskOutput(models.QueryTaskRequest{
		TaskID: &taskid,
	})
	if err != nil {
		return "", fmt.Errorf("get task output error: %v", err)
	}
	return output[taskid], nil
}

/*
实时输出
./bin/goosefs fsadmin report
*/
func Report() (string, error) {
	taskid, err := addTask(TaskRequest{
		Command: *config.Config.Bin,
		Args:    []string{"fsadmin", "report"},
	})
	if err != nil {
		return "", err
	}
	// log.Infof("taskid:", taskid)
	// wait for task done
	count := 0
	for {
		count++
		if count > 30 {
			return "", fmt.Errorf("wait for task done timeout, you can call output api to get task output, taskid: %s", taskid)
		}
		log.Debugf("get task status: %s", taskid)
		status, err := GetTaskStatus(models.FilterGoosefsTaskRequest{
			TaskID: &taskid,
		})
		if err != nil {
			return "", err
		}
		log.Debugf("task status: %s", status.Status)
		if status.Status == "<nil>" || status.Status == models.TaskStatusRunning {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	// 读取输出文件
	output, err := GetTaskOutput(models.QueryTaskRequest{
		TaskID:   &taskid,
		TaskName: nil,
	})
	if err != nil {
		return "", fmt.Errorf("get task output error: %v", err)
	}
	return output[taskid], nil
}

func GetCmdStatus(exitCode string) models.TaskState {
	switch exitCode {
	case "exit status 0":
		return models.TaskStatusSuccess
	case "<nil>":
		return models.TaskStatusRunning
	default:
		return models.TaskStatusFailed
	}
}
