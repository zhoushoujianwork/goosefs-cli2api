package api

import (
	"goosefs-cli2api/internal/executor"
	"goosefs-cli2api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SparkSubmit(c *gin.Context) {
	var req models.SparkSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := executor.SparkSubmit(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
