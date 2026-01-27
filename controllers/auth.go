package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type AuthController struct{}

var jwtSecret = []byte("your-secret-key")

type LoginRequest struct {
    Email    string `json:"email" form:"email"`
    Password string `json:"password" form:"password"`
}

type Claims struct {
    UserID int    `json:"user_id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// ═══════════════════════════════════════════
// WEB ROUTES (for browser)
// ═══════════════════════════════════════════

// GET /auth/login - Show login page
func (ac *AuthController) LoginPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login", gin.H{
        "title": "Admin Login",
    })
}

// POST /auth/login - Handle form submit
// POST /auth/login - Handle form submit
func (ac *AuthController) LoginSubmit(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.HTML(http.StatusOK, "login", gin.H{  // ← "login" not "login.html"
            "title": "Admin Login",
            "error": "Invalid input",
        })
        return
    }

    // Check credentials
    if req.Email == "admin@animeshop.com" && req.Password == "password123" {
        c.SetCookie("admin_logged_in", "true", 86400, "/", "", false, true)
        c.Redirect(http.StatusFound, "/admin/dashboard")
        return
    }

    c.HTML(http.StatusOK, "login", gin.H{  // ← HERE! Change "login.html" to "login"
        "title": "Admin Login",
        "error": "Invalid email or password",
    })
}


// ═══════════════════════════════════════════
// API ROUTES (for subsystems - JSON)
// ═══════════════════════════════════════════

// POST /api/auth/login - API login, returns JWT
func (ac *AuthController) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
        return
    }

    if req.Email == "admin@animeshop.com" && req.Password == "password123" {
        claims := &Claims{
            UserID: 1,
            Name:   "Admin User",
            Email:  req.Email,
            Role:   "super_admin",
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, _ := token.SignedString(jwtSecret)

        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "token":   tokenString,
            "user": gin.H{
                "id":    1,
                "name":  "Admin User",
                "email": req.Email,
                "role":  "super_admin",
            },
        })
        return
    }

    c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
}

// GET /api/auth/validate - Validate JWT token
func (ac *AuthController) Validate(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    if len(tokenString) > 7 {
        tokenString = tokenString[7:] // Remove "Bearer "
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
            "id":    claims.UserID,
            "name":  claims.Name,
            "email": claims.Email,
            "role":  claims.Role,
        },
    })
}
