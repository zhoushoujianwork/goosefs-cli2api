package api

import (
	"goosefs-cli2api/internal/executor"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	a := router.Group("api/v1")
	a.POST("/exec", executeTask)
	a.GET("/status/:task_id", getTaskStatus)
	a.GET("/output/:task_id", getTaskOutput)

	a.POST("/gfs", GoosefsExecute)
	a.GET("/gfs/report", GooseFSReport)
}

func executeTask(c *gin.Context) {
	var req executor.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID, err := executor.AddTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_id": taskID})
}

func getTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	status, err := executor.GetTaskStatus(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func getTaskOutput(c *gin.Context) {
	taskID := c.Param("task_id")
	output, err := executor.GetTaskOutput(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, output)
}
