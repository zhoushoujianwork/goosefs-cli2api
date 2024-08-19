package models

type QueryTaskRequest struct {
	TaskID   *string `json:"task_id"`
	TaskName *string `json:"task_name"`
}

type TaskState string

const (
	TaskStatusSuccess       TaskState = "success"
	TaskStatusFailed        TaskState = "failed"
	TaskStatusNotallSuccess TaskState = "notallsuccess"
	TaskStatusRunning       TaskState = "running"
)

// Data只输出有 OutPut内容的 格式如下{"taskid":"loadPath:记录数量"}
type TasksStatus struct {
	Data      map[string]TaskInfo `json:"data"` // 这里展示的是每个 Path任务的 taskID 和任务执行的 CMD结果
	Status    TaskState           `json:"status"`
	TotalTask int                 `json:"total_task"`
}

type TaskInfo struct {
	Path         string `json:"path"`          // 任务缓存的路径
	ExitCode     string `json:"exit_code"`     // exit_code 为 0 表示任务成功; <nil> 执行中
	SuccessCount int    `json:"success_count"` // 成功 load的对象数量
	Total        int    `json:"total"`         // 总共 load的对象数量
}
