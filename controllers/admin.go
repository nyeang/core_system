package controllers

import (
	"net/http"
	"time"

	"core-anime/config"
	"core-anime/models"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (ac *AdminController) Dashboard(c *gin.Context) {
	var totalUsers   int64
	var activeUsers  int64
	var todayLogins  int64
	var failedLogins int64

	today     := time.Now().Format("2006-01-02")
	last7Days := time.Now().AddDate(0, 0, -7)

	// Total subsystem users (exclude admins)
	config.DB.Model(&models.User{}).
		Where("role = ?", "user").
		Count(&totalUsers)

	// Active users = unique users logged in last 7 days
	config.DB.Model(&models.AuthLog{}).
		Where("created_at >= ? AND status = ? AND action = ?", last7Days, "success", "login").
		Distinct("user_id").
		Count(&activeUsers)

	// Today's successful logins
	config.DB.Model(&models.AuthLog{}).
		Where("DATE(created_at) = ? AND status = ? AND action = ?", today, "success", "login").
		Count(&todayLogins)

	// Today's failed logins
	config.DB.Model(&models.AuthLog{}).
		Where("DATE(created_at) = ? AND status = ? AND action = ?", today, "failed", "login").
		Count(&failedLogins)

	var recentLogs []models.AuthLog
	config.DB.Preload("User").Order("created_at desc").Limit(5).Find(&recentLogs)

	c.HTML(http.StatusOK, "admin", gin.H{
		"title": "Dashboard",
		"page":  "dashboard",
		"stats": gin.H{
			"totalUsers":   totalUsers,
			"activeUsers":  activeUsers,
			"todayLogins":  todayLogins,
			"failedLogins": failedLogins,
		},
		"recentLogs": recentLogs,
	})
}

func (ac *AdminController) User(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)

	c.HTML(http.StatusOK, "admin", gin.H{
		"title":    "User",
		"page":     "user",
		"userList": users,
		"user": gin.H{
			"name": "Admin",
		},
	})
}

func (ac *AdminController) Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "admin", gin.H{
		"title": "Settings",
		"page":  "settings",
	})
}

func (ac *AdminController) Log(c *gin.Context) {
	var totalLogs    int64
	var successCount int64
	var failedCount  int64

	today := time.Now().Format("2006-01-02")

	config.DB.Model(&models.AuthLog{}).
		Where("DATE(created_at) = ?", today).
		Count(&totalLogs)

	config.DB.Model(&models.AuthLog{}).
		Where("DATE(created_at) = ? AND status = ?", today, "success").
		Count(&successCount)

	config.DB.Model(&models.AuthLog{}).
		Where("DATE(created_at) = ? AND status = ?", today, "failed").
		Count(&failedCount)

	var subsystemStats []struct {
		Subsystem string
		Count     int64
	}
	config.DB.Model(&models.AuthLog{}).
		Select("subsystem, count(*) as count").
		Where("DATE(created_at) = ?", today).
		Group("subsystem").
		Scan(&subsystemStats)

	var authLogs []models.AuthLog
	config.DB.Preload("User").Order("created_at desc").Limit(50).Find(&authLogs)

	c.HTML(http.StatusOK, "admin", gin.H{
		"title": "Logs & Reports",
		"page":  "log",
		"stats": gin.H{
			"total":   totalLogs,
			"success": successCount,
			"failed":  failedCount,
		},
		"subsystemStats": subsystemStats,
		"log":            authLogs,
	})
}