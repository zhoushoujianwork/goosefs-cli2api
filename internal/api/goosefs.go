package api

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoosefsExecute(c *gin.Context) {
	var req models.GooseFSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch *req.Action {
	case models.GooseFSDistributeLoad:
		if req.Path == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}
		taskID, err := executor.DistrubuteLoad(*req.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task_id": taskID})
	case models.GooseFSLoadMetadata:
		if req.Path == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}
		taskID, err := executor.LoadMetadata(*req.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task_id": taskID})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "action not found"})
		return
	}
}

func GooseFSReport(c *gin.Context) {
	taskID, err := executor.Report()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task_id": taskID})
}
