package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type AdminController struct{}

func (ac *AdminController) Dashboard(c *gin.Context) {
    c.HTML(http.StatusOK, "dashboard.html", gin.H{
        "title": "Dashboard",
        "user": gin.H{
            "name": "Admin User",
            "email": "admin@example.com",
        },
        "stats": gin.H{
            "totalUsers": 1250,
            "totalOrders": 450,
            "revenue": "8,500",
            "newCustomers": 125,
        },
    })
}

func (ac *AdminController) Users(c *gin.Context) {
    c.HTML(http.StatusOK, "user.html", gin.H{
        "title": "Users Management",
        "user": gin.H{
            "name": "Admin User",
            "email": "admin@example.com",
        },
        "users": []gin.H{
            {"id": 1, "name": "John Doe", "email": "john@example.com", "role": "User"},
            {"id": 2, "name": "Jane Smith", "email": "jane@example.com", "role": "Admin"},
            {"id": 3, "name": "Bob Johnson", "email": "bob@example.com", "role": "User"},
        },
    })
}

func (ac *AdminController) Settings(c *gin.Context) {
    c.HTML(http.StatusOK, "settings.html", gin.H{
        "title": "Settings",
        "user": gin.H{
            "name": "Admin User",
            "email": "admin@example.com",
        },
    })
}