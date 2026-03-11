package controllers

import (
    "net/http"

    "core-anime/config"
    "core-anime/models"

    "github.com/gin-gonic/gin"
)

type LogController struct{}

func (lc *LogController) Create(c *gin.Context) {
    var req struct {
        UserID    uint   `json:"user_id"`
        Action    string `json:"action" binding:"required"`
        Subsystem string `json:"subsystem" binding:"required"`
        Details   string `json:"details"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
        return
    }

    log := models.AuthLog{
        UserID:    req.UserID,
        Action:    req.Action,
        IPAddress: c.ClientIP(),
        UserAgent: c.Request.UserAgent(),
        Subsystem: req.Subsystem,
        Status:    "success",
        Details:   req.Details,
    }

    if err := config.DB.Create(&log).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save log"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true})
}

func (lc *LogController) GetAll(c *gin.Context) {
    var logs []models.AuthLog
    config.DB.Preload("User").Order("created_at desc").Limit(100).Find(&logs)
    c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
}