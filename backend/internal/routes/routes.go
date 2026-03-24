// femabras/backend/internal/routes/routes.go
package routes

import (
	"net/http"
	"time"

	"github.com/Femabras/femabras/backend/internal/auth/handler"
	"github.com/Femabras/femabras/backend/internal/auth/service"
	"github.com/Femabras/femabras/backend/internal/config"
	"github.com/Femabras/femabras/backend/internal/handlers"
	"github.com/Femabras/femabras/backend/internal/middleware"

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

	// Public Webhook for Ad Network
	r.GET("/webhooks/ad-reward", handlers.AdRewardWebhook)

	// Auth Routes (Now using the new AuthHandler 'ah')
	r.POST("/register", ah.Register)
	r.POST("/verify-otp", ah.VerifyOTP)
	r.POST("/login", ah.Login)

	// Google Auth
	r.GET("/auth/google/login", ah.GoogleLogin)
	r.GET("/auth/google/callback", ah.GoogleCallback)

	// --- Protected Routes ---
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg, db))

	// Rate limit: 5 guesses per minute per IP
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	protected.Use(mgin.NewMiddleware(rateLimiter))

	protected.POST("/guess", ch.SubmitGuess)

	protected.GET("/challenge/attempts", ch.GetAttempts)
}
