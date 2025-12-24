package controllers

import (
	"net/http"

	"core-anime/config"
	"core-anime/models"

	"github.com/gin-gonic/gin"
)

// Exported function (capitalized)
func GetEpisode(c *gin.Context) {
	var episodes []models.Episode
	result := config.DB.Find(&episodes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, episodes)
}
