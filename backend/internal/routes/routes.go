// femabras/backend/internal/routes/routes.go
package routes

import (
	"net/http"
	"time"

	authHandler "github.com/Femabras/femabras/internal/auth/handler"
	authService "github.com/Femabras/femabras/internal/auth/service"
	challengeHandler "github.com/Femabras/femabras/internal/challenge/handler"
	challengeService "github.com/Femabras/femabras/internal/challenge/service"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/handlers"
	"github.com/Femabras/femabras/internal/middleware"
	"github.com/Femabras/femabras/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
	"gorm.io/gorm"
)

func Setup(
	r *gin.Engine,
	db *gorm.DB,
	cfg *config.Config,
	authSvc authService.AuthService,
	challengeSvc challengeService.ChallengeService,
) {
	ah := authHandler.NewAuthHandler(authSvc, cfg)

	// Single consolidated handler — flat handlers/challenge.go and
	// handlers/guess.go are now deleted. All challenge logic lives here.
	ch := challengeHandler.NewChallengeHandler(challengeSvc, db, cfg)

	redisClient := services.GetRedisClient()
	if redisClient == nil {
		panic("Redis is required for production rate limiting")
	}

	store, err := redisStore.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix:   "rate_limit:",
		MaxRetry: 3,
	})
	if err != nil {
		panic("failed to initialize redis rate limiter: " + err.Error())
	}

	// ── Public routes ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/challenge", ch.GetDailyChallenge)
	r.GET("/challenge/stream", ch.StreamStatus)
	r.GET("/webhooks/ad-reward", handlers.AdRewardWebhook(cfg))
	r.POST("/contact", handlers.ContactHandler(cfg))

	// ── Auth routes (5 req/min) ───────────────────────────────────────────────
	authRate := limiter.Rate{Period: 1 * time.Minute, Limit: 5}
	authLimiter := limiter.New(store, authRate, limiter.WithTrustForwardHeader(true))

	authGroup := r.Group("/")
	authGroup.Use(mgin.NewMiddleware(authLimiter))
	{
		authGroup.POST("/register", ah.Register)
		authGroup.POST("/verify-otp", ah.VerifyOTP)
		authGroup.POST("/login", ah.Login)
		authGroup.POST("/logout", ah.Logout)
		authGroup.POST("/refresh", ah.Refresh)
		authGroup.GET("/auth/google/login", ah.GoogleLogin)
		authGroup.GET("/auth/google/callback", ah.GoogleCallback)
	}

	// ── Protected routes (JWT + Redis cache) ─────────────────────────────────
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg, db))

	guessRate := limiter.Rate{Period: 1 * time.Minute, Limit: 10}
	guessLimiter := limiter.New(store, guessRate, limiter.WithTrustForwardHeader(true))

	protected.Use(mgin.NewMiddleware(guessLimiter))
	protected.POST("/guess", ch.SubmitGuess)
	protected.GET("/challenge/attempts", ch.GetAttempts)
	protected.POST("/claim", ch.ClaimPrize)
	protected.GET("/challenge/my-status", ch.GetMyStatus)
}
