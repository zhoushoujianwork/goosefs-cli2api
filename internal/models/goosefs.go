package models

const (
	GooseFSDistributeLoad = iota
	GooseFSLoadMetadata
)

type GooseFSRequest struct {
	Action *int    `json:"action" binding:"required"` // 必填 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata
	Path   *string `json:"path"`                      // 当 action 为 GooseFSDistributeLoad/GooseFSLoadMetadata 时必填
}
