package api

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"net/http"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
)

// @summary execute for goosefs cli
// @description 执行内置的 goosefs 命令，包括 distribute_load/load_metadata，返回 task_id，可以通过 task_id 获取执行状态或者输出
// @Tags GooseFS
// @Accept json
// @Produce json
// @Param req body models.GooseFSRequest true "DistrubuteLoad"
// @Success 200 {string} string
// @Router /api/v1/gfs [post]
func GoosefsExecute(c *gin.Context) {
	var req models.GooseFSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch *req.Action {
	case models.GooseFSDistributeLoad:
		if req.Path == nil {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		taskID, err := executor.DistrubuteLoad(*req.Path)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, taskID)
	case models.GooseFSLoadMetadata:
		if req.Path == nil {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		taskID, err := executor.LoadMetadata(*req.Path)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, taskID)
	case models.GooseFSList:
		if req.Path == nil {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		if req.TimeOut == nil {
			req.TimeOut = tea.Int(30) // 默认 30 秒
		}
		output, err := executor.List(*req.Path, *req.TimeOut)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, output)

	default:
		c.String(http.StatusBadRequest, "action not found, only support 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata 2: GooseFSList")
		return
	}
}

// @summary GooseFSReport
// @description GooseFSReport 获取 goosefs 集群状态
// @Tags GooseFS
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/gfs/report [get]
func GooseFSReport(c *gin.Context) {
	output, err := executor.Report()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
