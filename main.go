package main

import (
    "core-anime/config"  
    "core-anime/controllers"
    "github.com/gin-gonic/gin"
)

func main() {
    config.ConnectDatabase()
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

    authController := &controllers.AuthController{}
    adminController := &controllers.AdminController{}

    // API ROUTES
    api := r.Group("/api")
    api.POST("/auth/login", authController.Login)
    api.GET("/auth/validate", authController.Validate)

    // WEB ROUTES
    r.GET("/auth/login", authController.LoginPage)
    r.POST("/auth/login", authController.LoginSubmit)

    // Admin web routes
    r.GET("/admin/dashboard", adminController.Dashboard)
    r.GET("/admin/user", adminController.User)  
    r.GET("/admin/log", adminController.Log)   
    r.GET("/admin/settings", adminController.Settings)

    r.Run(":8080")
}
 