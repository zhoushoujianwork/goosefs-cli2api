package executor

import (
	"fmt"
	"goosefs-cli2api/config"
	"goosefs-cli2api/internal/models"
	"log"
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

/*
./bin/goosefs fs distributedLoad --replication 1 /data-datalake/deltalake/aaa.db/bbb/

支持多路径 Path 任务提交，返回 1 个 task_id
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
	// log.Println("taskid:", taskid)
	// wait for task done
	count := 0
	for {
		count++
		if count > timeOut {
			return "", fmt.Errorf("wait for task done timeout, you can call output api to get task output, taskid: %s", taskid)
		}
		// log.Println("get task status:", taskid)
		status, err := GetTaskStatus(models.QueryTaskRequest{
			TaskID: &taskid,
		})
		if err != nil {
			return "", err
		}
		log.Printf("task %s status: %s\n", taskid, status.Status)
		if status.Status == "<nil>" {
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
	// log.Println("taskid:", taskid)
	// wait for task done
	count := 0
	for {
		count++
		if count > 30 {
			return "", fmt.Errorf("wait for task done timeout, you can call output api to get task output, taskid: %s", taskid)
		}
		log.Println("get task status:", taskid)
		status, err := GetTaskStatus(models.QueryTaskRequest{
			TaskID: &taskid,
		})
		if err != nil {
			return "", err
		}
		log.Println("task status:", status.Status)
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
