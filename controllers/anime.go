package controllers

import (
    "net/http"
    "time"

    "core-anime/config"
    "core-anime/models"

    "github.com/gin-gonic/gin"
)

func GetAnimes(c *gin.Context) {
    var animes []models.Anime
    result := config.DB.Find(&animes)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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

// Admin page — fetch from local DB
func GetAnimeAdmin(c *gin.Context) {
    var animes []models.Anime
    config.DB.Find(&animes)

    c.HTML(http.StatusOK, "admin", gin.H{
        "page":  "anime",
        "anime": animes,
    })
}

func AddAnime(c *gin.Context) {
    var input struct {
        Title             string `json:"title"`
        Description       string `json:"description"`
        ReleaseDate       string `json:"release_date"`
        Genres            string `json:"genres"`
        ImageURL          string `json:"image_url"`
        SmallImageURL     string `json:"small_image_url"`
        LargeImageURL     string `json:"large_image_url"`
        ImageURLWebP      string `json:"image_url_webp"`
        SmallImageURLWebP string `json:"small_image_url_webp"`
        LargeImageURLWebP string `json:"large_image_url_webp"`
        TrailerURL        string `json:"trailer_url"`
        TrailerEmbedURL   string `json:"trailer_embed_url"`
        TrailerYoutubeID  string `json:"trailer_youtube_id"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
        return
    }

    var existing models.Anime
    if err := config.DB.Where("title = ?", input.Title).First(&existing).Error; err == nil {
        c.JSON(http.StatusOK, gin.H{"success": false, "message": "Already added"})
        return
    }

    releaseDate, _ := time.Parse(time.RFC3339, input.ReleaseDate)

    anime := models.Anime{
        Title:             input.Title,
        Description:       input.Description,
        ReleaseDate:       releaseDate,
        Genres:            input.Genres,
        ImageURL:          input.ImageURL,
        SmallImageURL:     input.SmallImageURL,
        LargeImageURL:     input.LargeImageURL,
        ImageURLWebP:      input.ImageURLWebP,
        SmallImageURLWebP: input.SmallImageURLWebP,
        LargeImageURLWebP: input.LargeImageURLWebP,
        TrailerURL:        input.TrailerURL,
        TrailerEmbedURL:   input.TrailerEmbedURL,
        TrailerYoutubeID:  input.TrailerYoutubeID,
    }

    if err := config.DB.Create(&anime).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to add"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true})
}