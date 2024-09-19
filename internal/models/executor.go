package models

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

type QueryTaskRequest struct {
	TaskID   *string `json:"task_id"`
	TaskName *string `json:"task_name"`
}

type TaskState string

func (t TaskState) TString() *string {
	return tea.String(string(t))
}

const (
	TaskStatusSuccess       TaskState = "success"
	TaskStatusFailed        TaskState = "failed"
	TaskStatusNotallSuccess TaskState = "notallsuccess"
	TaskStatusRunning       TaskState = "running"
	TaskStatusRestarted     TaskState = "restarted" //任务中断被系统重新执行
)

// Data只输出有 OutPut内容的 格式如下{"taskid":"loadPath:记录数量"}
type TasksStatus struct {
	Data      map[string]TaskInfo `json:"data"` // 这里展示的是每个 Path任务的 taskID 和任务执行的 CMD结果
	Status    TaskState           `json:"status"`
	TotalTask int                 `json:"total_task"`
}

type TaskInfo struct {
	TaskName     string    `json:"task_name"`     // 任务名称
	Path         string    `json:"path"`          // 任务缓存的路径
	ExitCode     string    `json:"exit_code"`     // exit_code 为 0 表示任务成功; <nil> 执行中
	SuccessCount int       `json:"success_count"` // 成功 load的对象数量
	TotalFile    int       `json:"total_file"`    // 总共 load的对象数量
	CreatedAt    time.Time `json:"created_at"`    // 任务创建时间 localtime
}
