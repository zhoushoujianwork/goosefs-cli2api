package models

type QueryTaskRequest struct {
	TaskID   *string `json:"task_id"`
	TaskName *string `json:"task_name"`
}

type TaskState string

const (
	TaskStatusSuccess TaskState = "success"
	TaskStatusFailed  TaskState = "failed"
	TaskStatusRunning TaskState = "running"
)

type TaskStatus struct {
	Data   map[string]string `json:"data"`
	Status TaskState         `json:"status"`
}
