package models

type GFSAction string

const (
	GooseFSDistributeLoad GFSAction = "GooseFSDistributeLoad"
	GooseFSLoadMetadata   GFSAction = "GooseFSLoadMetadata"
	GooseFSList           GFSAction = "GooseFSList"
)

type GooseFSRequest struct {
	Action   GFSAction `json:"action" binding:"required"` // 必填 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata 2: GooseFSList
	TaskName *string   `json:"task_name"`                 // 选填，支持提交多个任务到同一个任务标签上
	Path     []*string `json:"path" binding:"required"`   // 当 action 为 GooseFSDistributeLoad/GooseFSLoadMetadata/GooseFSList 时必填
	TimeOut  *int      `json:"timeout"`                   // 当 action 为 GooseFSList 由于没有挂起任务，所以需要指定超时时间 默认 30 秒
}
