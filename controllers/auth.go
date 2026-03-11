package controllers

import (
	"net/http"
	"time"

	"core-anime/config"
	"core-anime/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

var jwtSecret = []byte("your-secret-key")

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ═══════════════════════════════════════════
// WEB ROUTES (for browser / admin)
// ═══════════════════════════════════════════

func (ac *AuthController) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Admin Login",
	})
}

func (ac *AuthController) LoginSubmit(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "Admin Login",
			"error": "Invalid input",
		})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password",
		})
		return
	}

	if user.Role != "admin" && user.Role != "super_admin" {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "Admin Login",
			"error": "Access denied",
		})
		return
	}

	c.SetCookie("admin_logged_in", "true", 86400, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// ═══════════════════════════════════════════
// API ROUTES (for subsystems - JSON)
// ═══════════════════════════════════════════

// POST /api/auth/register
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	subsystem := c.GetHeader("X-Subsystem")
	if subsystem == "" {
		subsystem = "unknown"
	}

	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		config.DB.Create(&models.AuthLog{
			Action:    "register",
			IPAddress: c.ClientIP(),
			Subsystem: subsystem,
			Status:    "failed",
			Details:   "Email already registered: " + req.Email,
		})
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already registered"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}

	config.DB.Create(&models.AuthLog{
		UserID:    user.ID,
		Action:    "register",
		IPAddress: c.ClientIP(),
		Subsystem: subsystem,
		Status:    "success",
		Details:   "New user registered: " + user.Email,
	})

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// POST /api/auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	subsystem := c.GetHeader("X-Subsystem")
	if subsystem == "" {
		subsystem = "unknown"
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		config.DB.Create(&models.AuthLog{
			Action:    "login",
			IPAddress: c.ClientIP(),
			Subsystem: subsystem,
			Status:    "failed",
			Details:   "User not found: " + req.Email,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		config.DB.Create(&models.AuthLog{
			UserID:    user.ID,
			Action:    "login",
			IPAddress: c.ClientIP(),
			Subsystem: subsystem,
			Status:    "failed",
			Details:   "Wrong password for: " + req.Email,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
		return
	}

	config.DB.Create(&models.AuthLog{
		UserID:    user.ID,
		Action:    "login",
		IPAddress: c.ClientIP(),
		Subsystem: subsystem,
		Status:    "success",
		Details:   "Login from " + c.ClientIP(),
	})

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// GET /api/auth/me
func (ac *AuthController) Me(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 {
		tokenString = tokenString[7:]
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// GET /api/auth/validate
func (ac *AuthController) Validate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 {
		tokenString = tokenString[7:]
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       claims.UserID,
			"username": claims.Name,
			"email":    claims.Email,
			"role":     claims.Role,
		},
	})
}

// ═══════════════════════════════════════════
// Helper
// ═══════════════════════════════════════════

func generateToken(user models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Username,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}