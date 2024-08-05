package api

import (
	"goosefs-cli2api/internal/executor"
	"net/http"

	_ "goosefs-cli2api/docs"

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

	a.GET("/status/:task_id", getTaskStatus)
	a.GET("/output/:task_id", getTaskOutput)

	a.POST("/gfs", GoosefsExecute)
	a.GET("/gfs/report", GooseFSReport)
}

// @Summary GetTaskStatus
// @Description GetTaskStatus
// @Accept json
// @Produce json
// @Param task_id path string true "task_id"
// @Success 200 {object} executor.TaskStatus
// @Router /api/v1/status/{task_id} [get]
func getTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	status, err := executor.GetTaskStatus(taskID)
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
// @Param task_id path string true "task_id"
// @Success 200 {string} string
// @Router /api/v1/output/{task_id} [get]
func getTaskOutput(c *gin.Context) {
	taskID := c.Param("task_id")
	output, err := executor.GetTaskOutput(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, output)
}
