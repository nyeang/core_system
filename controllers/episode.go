package controllers

import (
    "net/http"
    "core-anime/config"
    "core-anime/models"
    "github.com/gin-gonic/gin"
)

func GetEpisode(c *gin.Context) {
    var episodes []models.Episode
    result := config.DB.Find(&episodes)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, episodes)
}

// ← ADD THIS NEW FUNCTION
func GetEpisodeByAnimeID(c *gin.Context) {
    animeID := c.Param("id")
    var episodes []models.Episode
    result := config.DB.Where("anime_id = ?", animeID).Find(&episodes)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    c.JSON(http.StatusOK, episodes)
}
