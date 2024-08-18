package api

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"net/http"

	"github.com/xops-infra/noop/log"

	_ "goosefs-cli2api/docs"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine) {
	a := router.Group("api/v1")

	// add swagger
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.Next()
	}, ginSwagger.WrapHandler(swaggerFiles.Handler))

	a.GET("/status", getTaskStatus)
	a.GET("/output", getTaskOutput)

	a.POST("/gfs", GoosefsExecute)
	a.GET("/gfs/report", GooseFSReport)
}

// @Summary GetTaskStatus
// @Description GetTaskStatus
// @Accept json
// @Produce json
// @Param task_id query string true "task_id"
// @Param task_name query string false "task_name"
// @Param action query string false "action"
// @Param status query string false "status"
// @Success 200 {object} models.QueryTaskRequest
// @Router /api/v1/status [get]
func getTaskStatus(c *gin.Context) {
	var req models.FilterGoosefsTaskRequest

	taskID := c.Query("task_id")
	if taskID != "" {
		req.TaskID = &taskID
	}
	taskName := c.Query("task_name")
	if taskName != "" {
		req.TaskName = &taskName
	}
	action := c.Query("action")
	if action != "" {
		gfsAction := models.GooseFSAction(action)
		req.Action = &gfsAction
	}
	taskStatus := c.Query("status")
	if taskStatus != "" {
		gfsStatus := models.TaskState(taskStatus)
		req.Status = &gfsStatus
	}
	log.Debugf(tea.Prettify(req))
	status, err := executor.GetTaskStatus(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// @Summary GetTaskOutput
// @Description GetTaskOutput
// @Accept json
// @Produce json
// @Param task_id query string true "task_id"
// @Param task_name query string false "task_name"
// @Success 200 {object} map[string]string
// @Router /api/v1/output [get]
func getTaskOutput(c *gin.Context) {
	var req models.QueryTaskRequest
	taskID := c.Query("task_id")
	if taskID != "" {
		req.TaskID = &taskID
	}
	taskName := c.Query("task_name")
	if taskName != "" {
		req.TaskName = &taskName
	}
	log.Debugf(tea.Prettify(req))
	if req.TaskID == nil && req.TaskName == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id or task_name is required"})
		return
	}
	output, err := executor.GetTaskOutput(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
