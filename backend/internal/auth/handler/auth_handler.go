// femabras/backend/internal/auth/handler/auth_handler.go
package handler

import (
	"crypto/rand"
	"crypto/subtle"
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
	return &AuthHandler{service: svc, cfg: cfg}
}

func (h *AuthHandler) isProd() bool {
	return h.cfg.FrontendURL != "http://localhost:3000"
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	isProd := h.isProd()
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
	isProd := h.isProd()
	domain := ""
	sameSite := http.SameSiteLaxMode
	if isProd {
		domain = ".femabras.com"
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "access_token", Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: isProd, Domain: domain, SameSite: sameSite,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "refresh_token", Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: isProd, Domain: domain, SameSite: sameSite,
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
	if n, err := fmt.Sscanf(req.UserID, "%d", &uid); err != nil || n != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	access, refresh, err := h.service.VerifyOTP(c.Request.Context(), uid, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookies(c, access, refresh)
	c.JSON(http.StatusOK, gin.H{"message": "Verification successful"})
}

// GoogleLogin generates a cryptographically random CSRF state token, stores it
// in a short-lived HttpOnly cookie, then redirects the user to Google.
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state token"})
		return
	}
	state := fmt.Sprintf("%x", b)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes — enough for the OAuth round-trip
		HttpOnly: true,
		Secure:   h.isProd(),
		SameSite: http.SameSiteLaxMode, // Lax allows the cookie on the return redirect from Google
	})

	url := service.GetGoogleAuthURL(state)
	if url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback verifies the CSRF state, exchanges the code for tokens, sets
// HttpOnly cookies, and redirects to the frontend — no tokens in the URL.
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing OAuth state cookie"})
		return
	}

	// Use constant-time comparison to prevent timing oracle on state tokens
	if subtle.ConstantTimeCompare([]byte(c.Query("state")), []byte(stateCookie)) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	// Immediately clear the one-time state cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "oauth_state", Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
	})

	code := c.Query("code")
	access, refresh, err := h.service.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth failed"})
		return
	}

	h.setAuthCookies(c, access, refresh)
	// Redirect to the frontend root — no tokens in the URL
	c.Redirect(http.StatusFound, h.cfg.FrontendURL)
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
		h.clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookies(c, access, newRefresh)
	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}
