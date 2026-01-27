package routes

import (
	"core-anime/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		api.GET("/users", controllers.GetUser)
		api.GET("/users/:id", controllers.GetUserByID)
	}

	
}

func SetupRoutes(r *gin.Engine) {
    adminCtrl := &controllers.AdminController{}
    
    // Admin routes
    admin := r.Group("/admin")
    {
        admin.GET("/dashboard", adminCtrl.Dashboard)
        admin.GET("/user", adminCtrl.User)
        admin.GET("/logs", adminCtrl.Log)
        admin.GET("/settings", adminCtrl.Settings)
    }
    
    // Redirect root to dashboard
    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
    })
}



