package executor

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/models"
	"goosefs-cli2api/pkg/dingtalk"
	"time"

	"github.com/xops-infra/noop/log"

	"github.com/alibabacloud-go/tea/tea"
)

// 实现任务完成后的告警通知
func checkTasksIsFinished(act models.GFSAction, taskids []string, task_name *string) {

	for {
		time.Sleep(5 * time.Second)
		log.Debugf("checkTasksIsFinished")
		var status []models.TaskStatus
		if task_name != nil && *task_name != "" {
			// 按照任务进度查询
			_status, err := GetTaskStatus(models.QueryTaskRequest{
				TaskName: task_name,
			})
			if err != nil {
				log.Errorf("checkTasksIsFinished error: %s", err)
				continue
			}
			status = append(status, _status)
		} else {
			// 按照任务ID查询
			for _, id := range taskids {
				_status, err := GetTaskStatus(models.QueryTaskRequest{
					TaskID: tea.String(id),
				})
				if err != nil {
					log.Errorf("checkTasksIsFinished error: %s", err)
					continue
				}
				status = append(status, _status)
			}
		}

		allStatus := make(map[models.TaskState]int, 0)
		for _, s := range status {
			if s.Status == models.TaskStatusRunning {
				allStatus[models.TaskStatusRunning]++
				continue
			}
			allStatus[s.Status]++
		}

		var msg string
		if allStatus[models.TaskStatusFailed] > 0 {
			msg = "告警:" + string(act) + "\tERROR!\n"
		} else {
			msg = "通知:" + string(act) + "\tSUCCESS!\n"
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
			Name:    tea.StringValue(req.TaskName),
			Command: *config.Config.Bin,
			Args:    []string{"fs", "distributedLoad", "--replication", "1", *p},
		})
		if err != nil {
			return nil, err
		}
		taskids = append(taskids, taskID)
	}
	go checkTasksIsFinished(models.GooseFSDistributeLoad, taskids, req.TaskName)
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
			Name:    tea.StringValue(req.TaskName),
			Command: *config.Config.Bin,
			Args:    []string{"fs", "loadMetadata", "-R", *p},
		})
		if err != nil {
			return nil, err
		}
		taskids = append(taskids, taskID)
	}
	go checkTasksIsFinished(models.GooseFSLoadMetadata, taskids, req.TaskName)
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
		status, err := GetTaskStatus(models.QueryTaskRequest{
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
		status, err := GetTaskStatus(models.QueryTaskRequest{
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
