package handlers

import (
    "net/http"
    "todo-list-api/internal/services"
    "github.com/gin-gonic/gin"
)

// RegisterInput defines the input for user registration
type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterInput true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func Register(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input RegisterInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := authService.Register(input.Username, input.Password); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
    }
}

// LoginInput defines the input for user login
type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// @Summary Login a user
// @Description Login a user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginInput true "User data"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input LoginInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        accessToken, refreshToken, err := authService.Login(input.Username, input.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "access_token":  accessToken,
            "refresh_token": refreshToken,
        })
    }
}

// RefreshInput defines the input for token refresh
type RefreshInput struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshInput true "Refresh token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func Refresh(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input RefreshInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        newAccessToken, err := authService.RefreshToken(input.RefreshToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
    }
}