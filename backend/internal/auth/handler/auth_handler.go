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

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken string, refreshToken string) {
	isProd := h.cfg.FrontendURL != "http://localhost:3000"
	domain := ""
	sameSite := http.SameSiteLaxMode

	if isProd {
		domain = ".femabras.com"
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   900,
		HttpOnly: true,
		Secure:   isProd,
		Domain:   domain,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   604800,
		HttpOnly: true,
		Secure:   isProd,
		Domain:   domain,
		SameSite: sameSite,
	})
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	isProd := h.cfg.FrontendURL != "http://localhost:3000"
	domain := ""
	sameSite := http.SameSiteLaxMode

	if isProd {
		domain = ".femabras.com"
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "access_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: isProd, Domain: domain, SameSite: sameSite,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "refresh_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: isProd, Domain: domain, SameSite: sameSite,
	})
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
		"message": "User created. OTP sent to email.",
		"user_id": userID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	access, refresh, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookies(c, access, refresh)
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

	access, refresh, err := h.service.VerifyOTP(c.Request.Context(), uid, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookies(c, access, refresh)
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
	access, refresh, err := h.service.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth failed"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth/success?access=%s&refresh=%s", h.cfg.FrontendURL, access, refresh))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil && refreshToken != "" {
		h.service.RevokeRefreshToken(c.Request.Context(), refreshToken)
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	access, newRefresh, err := h.service.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearAuthCookies(c) // Wipe dead cookies
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookies(c, access, newRefresh)
	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}
