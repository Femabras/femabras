// femabras/backend/internal/routes/routes.go
package routes

import (
	"net/http"
	"time"

	"github.com/Femabras/femabras/internal/auth/handler"
	"github.com/Femabras/femabras/internal/auth/service"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/handlers"
	"github.com/Femabras/femabras/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, cfg *config.Config, authSvc service.AuthService) {
	// Initialize Auth Dependencies for the routes
	ah := handler.NewAuthHandler(authSvc, cfg)

	// Initialize Other Handlers
	ch := &handlers.ChallengeHandler{DB: db}

	// --- Public Routes ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/challenge", ch.GetDailyChallenge)
	// Webhook for Ad Network
	r.GET("/webhooks/ad-reward", handlers.AdRewardWebhook)

	// --- 🟢 Auth Routes (STRICTLY THROTTLED to protect Twilio & Passwords) ---
	// 5 requests per minute per IP
	authRate := limiter.Rate{Period: 1 * time.Minute, Limit: 5}
	authStore := memory.NewStore()
	authLimiter := limiter.New(authStore, authRate, limiter.WithTrustForwardHeader(true))

	authGroup := r.Group("/")
	authGroup.Use(mgin.NewMiddleware(authLimiter))
	{
		authGroup.POST("/register", ah.Register)
		authGroup.POST("/verify-otp", ah.VerifyOTP)
		authGroup.POST("/login", ah.Login)
		authGroup.GET("/auth/google/login", ah.GoogleLogin)
		authGroup.GET("/auth/google/callback", ah.GoogleCallback)
	}

	// --- Protected Routes (Throttled to prevent guess spamming) ---
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg, db))

	guessRate := limiter.Rate{Period: 1 * time.Minute, Limit: 5}
	guessStore := memory.NewStore()
	guessLimiter := limiter.New(guessStore, guessRate, limiter.WithTrustForwardHeader(true))

	protected.Use(mgin.NewMiddleware(guessLimiter))
	protected.POST("/guess", ch.SubmitGuess)
	protected.GET("/challenge/attempts", ch.GetAttempts)
}
