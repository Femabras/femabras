// femabras/backend/internal/auth/handler/auth_handler.go
package handler

import (
	"fmt"
	"net/http"

	"github.com/Femabras/femabras/internal/auth/service"
	"github.com/Femabras/femabras/internal/auth/types"
	"github.com/Femabras/femabras/internal/config"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(svc service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: svc,
		cfg:     cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userID, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created. OTP sent to phone.",
		"user_id": userID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req types.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var uid uint
	fmt.Sscanf(req.UserID, "%d", &uid)

	token, err := h.service.VerifyOTP(c.Request.Context(), uid, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "auth_token", Value: token, Path: "/", MaxAge: 86400,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Verification successful"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := service.GetGoogleAuthURL("state-token")
	if url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "state-token" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := h.service.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth failed"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/success?token=%s", h.cfg.FrontendURL, token))
}
