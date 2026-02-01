package handler

import (
	"net/http"

	"geo/server/backend/service"

	"github.com/gin-gonic/gin"
)

func GetPendingTasks(c *gin.Context) {
	tasks, err := service.GetPendingTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func UpdateTaskResult(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTaskResultReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.UpdateTaskResult(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
