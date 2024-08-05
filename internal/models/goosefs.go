package models

const (
	GooseFSDistributeLoad = iota
	GooseFSLoadMetadata
)

type GooseFSRequest struct {
	Action *int    `json:"action" binding:"required"`
	Path   *string `json:"path"` // 当 action 为 GooseFSDistributeLoad/GooseFSLoadMetadata 时必填
}
