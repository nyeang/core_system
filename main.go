package main

import (
	"core-anime/config"
	"core-anime/models"
	"core-anime/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    config.ConnectDatabase()

    var count int64
    config.DB.Model(&models.Anime{}).Count(&count)
    if count == 0 {
        log.Println("Seeding anime from Jikan...")
        SeedAnimeFromJikan()
    }


    r := gin.Default()
    r.LoadHTMLGlob("templates/**/*.html")

    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    routes.SetupRoutes(r)
    r.Run(":8083")
}
