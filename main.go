package main

import (
	"core-anime/config"
	"core-anime/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*")
	
	r.Static("/static", "./static")

	config.ConnectDatabase()

	routes.RegisterRoutes(r)
	routes.SetupRoutes(r)

	r.Run(":8080")
}