package controllers

import (
    "net/http"

    "core-anime/config"
    "core-anime/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

func GetUser(c *gin.Context) {
    var users []models.User
    if err := config.DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
    var user models.User
    if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
        Role     string `json:"role"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check duplicate email
    var existing models.User
    if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    role := req.Role
    if role == "" {
        role = "user"
    }

    user := models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashed),
        Role:         role,
    }
    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"success": true, "data": user})
}

func DeleteUser(c *gin.Context) {
    var user models.User
    if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    config.DB.Delete(&user)
    c.JSON(http.StatusOK, gin.H{"success": true})
}