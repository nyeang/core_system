package controllers

import (
	"net/http"

	"core-anime/config"
	"core-anime/models"

	"github.com/gin-gonic/gin"
)

// Exported function (capitalized)
func GetAnimes(c *gin.Context) {
	var animes []models.Anime
	result := config.DB.Find(&animes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		c.IndentedJSON(http.StatusOK, animes)
		return
	}
	c.IndentedJSON(http.StatusOK, animes)
}


func GetAnimeByID(c *gin.Context) {
	id := c.Param("id")

	var anime models.Anime

	if err := config.DB.First(&anime, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	c.JSON(http.StatusOK, anime)
}