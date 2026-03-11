package routes

import (
    "core-anime/controllers"
    "net/http"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    adminCtrl := &controllers.AdminController{}
    authCtrl  := &controllers.AuthController{}
    logCtrl   := &controllers.LogController{}

    // Redirect root
    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
    })

    // Admin web routes
    admin := r.Group("/admin")
    {
        admin.GET("/dashboard", adminCtrl.Dashboard)
        admin.GET("/user",      adminCtrl.User)
        admin.GET("/logs",      adminCtrl.Log)
        admin.GET("/settings",  adminCtrl.Settings)
    }

    // API routes for subsystems
    api := r.Group("/api")
    {
        api.POST("/auth/login",        authCtrl.Login)
        api.POST("/auth/register",     authCtrl.Register)
        api.GET("/auth/validate",      authCtrl.Validate)
        api.GET("/auth/me",            authCtrl.Me)
        api.POST("/logs",              logCtrl.Store)
        api.GET("/anime",              controllers.GetAnimes)
        api.GET("/anime/:id",          controllers.GetAnimeByID)
        api.GET("/anime/:id/episodes", controllers.GetEpisodeByAnimeID)
        api.GET("/episodes",           controllers.GetEpisode)
        api.GET("/users",              controllers.GetUser)
        api.GET("/users/:id",          controllers.GetUserByID)
    }
}