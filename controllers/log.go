package controllers

import (
	"fmt"
	"net/http"

	"core-anime/config"
	"core-anime/models"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

type CreateLogRequest struct {
	UserID    interface{} `json:"user_id"`
	Action    string      `json:"action"`
	Subsystem string      `json:"subsystem"`
	Details   string      `json:"details"`
}

func (lc *LogController) Store(c *gin.Context) {
	var req CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if req.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "action is required"})
		return
	}

	// Convert user_id to uint (handles both string and number)
	var userID uint
	switch v := req.UserID.(type) {
	case float64:
		userID = uint(v)
	case string:
		// Try to parse string to uint
		var parsed uint64
		_, err := fmt.Sscanf(v, "%d", &parsed)
		if err == nil {
			userID = uint(parsed)
		}
	}

	log := models.AuthLog{
		UserID:    userID,
		Action:    req.Action,
		Subsystem: req.Subsystem,
		Details:   req.Details,
		IPAddress: c.ClientIP(),
		Status:    "success",
	}

	if err := config.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save log"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true})
}