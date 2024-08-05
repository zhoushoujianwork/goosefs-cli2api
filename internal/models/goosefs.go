package models

const (
	GooseFSDistributeLoad = iota
	GooseFSLoadMetadata
	GooseFSList
)

type GooseFSRequest struct {
	Action  *int    `json:"action" binding:"required"` // 必填 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata 2: GooseFSList
	Path    *string `json:"path"`                      // 当 action 为 GooseFSDistributeLoad/GooseFSLoadMetadata/GooseFSList 时必填
	TimeOut *int    `json:"timeout"`                   // 当 action 为 GooseFSList 由于没有挂起任务，所以需要指定超时时间 默认 30 秒
}
