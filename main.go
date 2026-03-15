package main

import (
	"time"
	"core-anime/config"
	"core-anime/models"
	"core-anime/routes"
	"github.com/gin-contrib/cors"
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

    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://165.22.250.160",
            "http://188.166.184.64",
            "http://152.42.220.220",
        },
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Subsystem",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

    routes.SetupRoutes(r)
    r.Run(":8083")
}
