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

func checkTasksIsFinished(act models.GooseFSAction, taskName *string) {
	if taskName == nil {
		log.Debugf("task_name and taskids are empty, skip checkTasksIsFinished")
		return
	}

	for {
		time.Sleep(5 * time.Second)
		log.Debugf("checkTasksIsFinished")

		status, err := GetTaskStatus(models.FilterGoosefsTaskRequest{
			TaskName: taskName,
			Action:   &act,
		})
		if err != nil {
			log.Errorf("checkTasksIsFinished error: %s", err)
			break
		}

		var msg string
		switch status.Status {
		case models.TaskStatusRunning:
			log.Debugf("task %s is running wait for it done", tea.StringValue(taskName))
			continue
		case models.TaskStatusSuccess:
			msg = "通知:" + string(act) + " " + string(status.Status) + " for task " + tea.StringValue(taskName) + "\n"
		default:
			msg = "告警:" + string(act) + " " + string(status.Status) + " for task " + tea.StringValue(taskName) + "\n"
		}

		// 个性解析输出结果的数据到数据库

		output, err := GetTaskOutput(models.QueryTaskRequest{
			TaskName: taskName,
		})
		if err != nil {
			log.Errorf("GetTaskOutput error: %v", err)
			break
		}

		for taskId, out := range output {
			var successCount, loadedCount int
			defer func() {
				config.DB.UpdateGoosefsTask(taskId, models.UpdateGoosefsTaskRequest{
					SuccessCount: tea.Int(successCount),
					Total:        tea.Int(loadedCount + successCount),
				})
			}()
			lines := strings.Split(out, "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				switch act {
				case models.GFSDistributeLoad, models.GFSForceLoad:
					// 只统计Successfully loaded path的行数
					if strings.Contains(line, "Successfully loaded path") {
						successCount++
					} else if strings.Contains(line, "is already fully loaded in GooseFS") {
						loadedCount++
					}
				default:
					// 输出结果不需要统计成功失败的
					continue
				}
			}

			var taskPath string
			if _, ok := status.Data[taskId]; ok {
				taskPath = status.Data[taskId].Path
			} else {
				taskPath = taskId
			}
			if len(msg) > 20000-100 {
				msg += "..."
			} else {
				msg += fmt.Sprintf("task: %s %d/%d (success)/(total)\n", taskPath, successCount, loadedCount+successCount)
			}
		}

		// 优化降低通知文本 body 大小不合法;solution:请保持大小在 20000bytes
		// for _, taskInfo := range status.Data {
		// 	msg += fmt.Sprintf("path: %s %d/%d\n", taskInfo.Path, taskInfo.SuccessCount, taskInfo.TotalFile)
		// 	if len(msg) > 20000-5 {
		// 		msg += "..."
		// 		break
		// 	}
		// }
		dingtalk.SendAlert(msg)
		break
	}
}

/*
./bin/goosefs fs distributedLoad --replication 1 /data-datalake/deltalake/aaa.db/bbb/

支持多路径 Path 任务提交，返回 1 个 task_id
实现任务完成告警，每次任务进来后挂起一个任务查询的协程，直到任务状态不是 running；
*/

func DistrubuteLoad(req models.GooseFSRequest) (models.GooseFSExecuteResponse, error) {
	results := make([]models.Result, 0)
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return models.GooseFSExecuteResponse{}, fmt.Errorf("path is required, should not be empty")
		}

		taskID, err := addTask(TaskRequest{
			Action:   models.GFSDistributeLoad,
			Path:     *p,
			TaskName: tea.StringValue(req.TaskName),
			Command:  *config.Config.Bin,
			Args:     []string{"fs", "distributedLoad", "--replication", "1", *p},
		})
		if err != nil {
			return models.GooseFSExecuteResponse{}, err
		}

		// 创建 task 记录
		err = config.DB.CreateGoosefsTask(taskID, models.GoosefsTaskRequest{
			TaskName: req.TaskName,
			Path:     p,
			Action:   req.Action,
		})
		if err != nil {
			return models.GooseFSExecuteResponse{}, err
		}

		results = append(results, models.Result{
			TaskID: taskID,
			Path:   *p,
		})
	}
	go checkTasksIsFinished(models.GFSDistributeLoad, req.TaskName)
	log.Infof("add DistrubuteLoad success for req: %s", tea.Prettify(req))
	return models.GooseFSExecuteResponse{Results: results, Total: len(results)}, nil
}

/*
./bin/goosefs fs loadMetadata -R /data-datalake/deltalake/aaa.db/bbb/
*/
func LoadMetadata(req models.GooseFSRequest) (models.GooseFSExecuteResponse, error) {
	results := make([]models.Result, 0)
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return models.GooseFSExecuteResponse{}, fmt.Errorf("path is required, should not be empty")
		}
		taskID, err := addTask(TaskRequest{
			TaskName: tea.StringValue(req.TaskName),
			Command:  *config.Config.Bin,
			Args:     []string{"fs", "loadMetadata", "-R", "-F", *p},
			Path:     *p,
			Action:   models.GFSLoadMetadata,
		})
		if err != nil {
			return models.GooseFSExecuteResponse{}, err
		}
		// 创建 task 记录
		err = config.DB.CreateGoosefsTask(taskID, models.GoosefsTaskRequest{
			TaskName: req.TaskName,
			Path:     p,
			Action:   req.Action,
		})
		if err != nil {
			return models.GooseFSExecuteResponse{}, err
		}

		results = append(results, models.Result{
			TaskID: taskID,
			Path:   *p,
		})
	}
	go checkTasksIsFinished(models.GFSLoadMetadata, req.TaskName)
	log.Infof("add LoadMetadata success for req: %s", tea.Prettify(req))
	return models.GooseFSExecuteResponse{Results: results, Total: len(results)}, nil
}

/*
GooseFSForceLoad 该步骤执行的是先去 LoadMetadata，然后再去 DistributeLoad，这样彻底更新
先执行 LoadMetadata，成功后再执行 DistributeLoad
*/
func ForceLoad(req models.GooseFSRequest) error {
	for _, p := range req.Path {
		if p == nil || *p == "" {
			return fmt.Errorf("path is required, should not be empty")
		}
		// 防止阻塞，线程执行后续任务
		go forceLoadExector(p, req)
	}
	go checkTasksIsFinished(models.GFSForceLoad, req.TaskName)
	log.Infof("add ForceLoad success for req: %s", tea.Prettify(req))
	return nil
}

func forceLoadExector(p *string, req models.GooseFSRequest) {
	_, err := runCmd(*config.Config.Bin, []string{"fs", "loadMetadata", "-R", "-F", *p})
	if err != nil {
		log.Errorf("loadMetadata error: %s", err)
	}
	taskID, err := addTask(TaskRequest{
		TaskName: tea.StringValue(req.TaskName),
		Command:  *config.Config.Bin,
		Args:     []string{"fs", "distributedLoad", "--replication", "1", *p},
		Path:     *p,
		Action:   models.GFSForceLoad,
	})
	if err != nil {
		log.Errorf("addTask error: %s for path: %s", err, *p)
	}
	// 创建 task 记录
	err = config.DB.CreateGoosefsTask(taskID, models.GoosefsTaskRequest{
		TaskName: req.TaskName,
		Path:     p,
		Action:   req.Action,
	})
	if err != nil {
		log.Errorf("create task error: %s for path: %s", err, *p)
	}
	log.Infof("add success taskid: %s for path: %s", taskID, *p)
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
	case "restarted":
		return models.TaskStatusRestarted
	default:
		return models.TaskStatusFailed
	}
}
