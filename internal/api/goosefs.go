package api

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"net/http"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
)

// @summary execute for goosefs cli
// @description 执行内置的 goosefs 命令，包括 distribute_load/load_metadata，返回 task_id，可以通过 task_id 获取执行状态或者输出;注意GooseFSList是等待执行的不是挂起的任务，且只支持 1 个 path查询。
// @Tags GooseFS
// @Accept json
// @Produce json
// @Param req body models.GooseFSRequest true "DistrubuteLoad"
// @Success 200 {object} models.GooseFSExecuteResponse
// @Router /api/v1/gfs [post]
func GoosefsExecute(c *gin.Context) {
	var req models.GooseFSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	switch req.Action {
	case models.GFSForceLoad:
		if req.Path == nil {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		resp, err := executor.ForceLoad(req)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)

	case models.GFSDistributeLoad:
		if req.Path == nil {
			c.JSON(http.StatusBadRequest, "path is required")
			return
		}
		resp, err := executor.DistrubuteLoad(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	case models.GFSLoadMetadata:
		if req.Path == nil {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		resp, err := executor.LoadMetadata(req)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	case models.GFSList:

		if req.TimeOut == nil {
			req.TimeOut = tea.Int(30) // 默认 30 秒
		}
		if len(req.Path) != 1 {
			c.String(http.StatusBadRequest, "list only support one path, cause this func is not running in background")
			return
		}
		if req.Path == nil || *req.Path[0] == "" {
			c.String(http.StatusBadRequest, "path is required")
			return
		}
		output, err := executor.List(*req.Path[0], *req.TimeOut)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, output)

	default:
		c.String(http.StatusBadRequest, "action not found, only support 0: GooseFSDistributeLoad 1: GooseFSLoadMetadata 2: GooseFSList 3: GooseFSReport")
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
