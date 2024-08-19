package models

type GooseFSAction string

const (
	GFSForceLoad      GooseFSAction = "GooseFSForceLoad"      // 该步骤执行的是先去 LoadMetadata，然后再去 DistributeLoad，这样彻底更新
	GFSDistributeLoad GooseFSAction = "GooseFSDistributeLoad" // 缓存数据，他的依据是 Master 上的 metadata
	GFSLoadMetadata   GooseFSAction = "GooseFSLoadMetadata"   // 只更新元数据信息，可以更新掉cos上变更的内容
	GFSList           GooseFSAction = "GooseFSList"
)

// 外部请求支持多路径
type GooseFSRequest struct {
	Action   GooseFSAction `json:"action" binding:"required"` // 必填 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata 2: GooseFSList 3: GooseFSForceLoad
	TaskName *string       `json:"task_name"`                 // 选填，支持提交多个任务到同一个任务标签上
	Path     []*string     `json:"path" binding:"required"`   // 当 action 为 GooseFSDistributeLoad/GooseFSLoadMetadata/GooseFSList 时必填
	TimeOut  *int          `json:"timeout"`                   // 当 action 为 GooseFSList 由于没有挂起任务，所以需要指定超时时间 默认 30 秒
}

// 内部处理都是单一路径执行的 task
type GoosefsTaskRequest struct {
	Action   GooseFSAction `gorm:"column:action;type:text" json:"action"`
	TaskName *string       `gorm:"column:task_name;type:text" json:"task_name"`
	Path     *string       `gorm:"column:path;type:text" json:"path"`
}

func (r *GoosefsTaskRequest) ToGoosefsTask(taskID string) *GoosefsTask {
	return &GoosefsTask{
		ID:       taskID,
		Action:   r.Action,
		TaskName: r.TaskName,
		Path:     r.Path,
	}
}

type GoosefsTask struct {
	ID           string        `gorm:"column:id;type:text;not null;primary_key" json:"id"`
	Action       GooseFSAction `gorm:"column:action;type:text" json:"action"`
	TaskName     *string       `gorm:"column:task_name;type:text" json:"task_name"`
	Path         *string       `gorm:"column:path;type:text" json:"path"`
	ExitCode     *string       `gorm:"column:exit_code;type:text" json:"exit_code"`
	SuccessCount *int          `gorm:"column:success_count;type:text" json:"success_count"`
	Total        *int          `gorm:"column:total;type:text" json:"total"`
}

func (*GoosefsTask) TableName() string {
	return "goosefs_task"
}

type UpdateGoosefsTaskRequest struct {
	ExitCode     *string `json:"exit_code"`
	SuccessCount *int    `json:"success_count"`
	Total        *int    `json:"total"`
}

type FilterGoosefsTaskRequest struct {
	TaskID   *string        `json:"task_id"`
	TaskName *string        `json:"task_name"`
	Action   *GooseFSAction `json:"action"`
	Status   *TaskState     `json:"status"`
}

// type GoosefsPathStatus struct {
// 	ID       string  `gorm:"column:id;type:text;not null;primary_key"`
// 	TaskID   *string `gorm:"column:task_id;type:text;not null"`
// 	TaskName *string `gorm:"column:task_name;type:text"`
// 	Path     *string `gorm:"column:path;type:text;not null"`
// 	ExitCode *string `gorm:"column:exit_code;type:text"`
// 	Count    *int    `gorm:"column:count;type:text"`
// }

// func (*GoosefsPathStatus) TableName() string {
// 	return "goosefs_path_status"
// }

// type CreatePathStatusRequest struct {
// 	TaskName *string `json:"task_name"`
// 	Path     *string `json:"path" binding:"required"`
// 	TaskID   *string `json:"task_id" binding:"required"`
// }

// func (r *CreatePathStatusRequest) ToPathStatus() *GoosefsPathStatus {
// 	return &GoosefsPathStatus{
// 		ID:       uuid.New().String(),
// 		TaskID:   r.TaskID,
// 		TaskName: r.TaskName,
// 		Path:     r.Path,
// 		ExitCode: nil,
// 		Count:    nil,
// 	}
// }

// type UpdatePathStatusRequest struct {
// 	ExitCode *string `gorm:"column:exit_code;type:text"`
// 	Count    *int    `gorm:"column:count;type:text"`
// }

// type FilterPathStatusRequest struct {
// 	// ID       *string `json:"id"` // 单查这个没意义，用下面的参数查
// 	TaskName *string `json:"task_name"`
// 	TaskID   *string `json:"task_id"`
// 	Path     *string `json:"path"`
// }
