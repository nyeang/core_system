package routes

import (
    "core-anime/controllers"
    "net/http"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    adminCtrl := &controllers.AdminController{}
    authCtrl  := &controllers.AuthController{}
    logCtrl   := &controllers.LogController{}  // ← add this

    // Redirect root
    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
    })

    // Auth web routes
    r.GET("/auth/login",  authCtrl.LoginPage)
    r.POST("/auth/login", authCtrl.LoginSubmit)

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
        // Auth
        api.POST("/auth/register", authCtrl.Register)
        api.POST("/auth/login",    authCtrl.Login)
        api.GET("/auth/me",        authCtrl.Me)
        api.GET("/auth/validate",  authCtrl.Validate)

        // Anime
        api.GET("/anime",              controllers.GetAnimes)
        api.GET("/anime/:id",          controllers.GetAnimeByID)
        api.GET("/anime/:id/episodes", controllers.GetEpisodeByAnimeID)
        api.GET("/episodes",           controllers.GetEpisode)

        // Users
        api.GET("/users",     controllers.GetUser)
        api.GET("/users/:id", controllers.GetUserByID)

        // Logs
        api.POST("/logs", logCtrl.Create)
        api.GET("/logs",  logCtrl.GetAll)
    }
}