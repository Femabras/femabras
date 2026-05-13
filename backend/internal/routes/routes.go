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
	isProd := cfg.FrontendURL != "http://localhost:3000"

	// Security headers on every response — CSP, HSTS, frame deny, etc.
	r.Use(middleware.SecurityHeaders())

	ah := authHandler.NewAuthHandler(authSvc, cfg)
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

	// CSRF token issuance — frontend calls this on app boot to obtain the
	// double-submit cookie. Safe to call repeatedly; idempotent.
	r.GET("/csrf", middleware.IssueCSRFToken(isProd))

	r.GET("/challenge", ch.GetDailyChallenge)
	r.GET("/challenge/stream", ch.StreamStatus)
	r.GET("/webhooks/ad-reward", handlers.AdRewardWebhook(cfg))

	// Contact form — strict rate limit: 3 submissions per hour per IP.
	// Without this, anyone can spam the admin inbox or burn Resend quota.
	contactRate := limiter.Rate{Period: 1 * time.Hour, Limit: 3}
	contactLimiter := limiter.New(store, contactRate, limiter.WithTrustForwardHeader(true))
	r.POST("/contact",
		mgin.NewMiddleware(contactLimiter),
		middleware.CSRF(isProd),
		handlers.ContactHandler(cfg),
	)

	// ── Auth routes (5 req/min, CSRF enforced on POST) ────────────────────────
	authRate := limiter.Rate{Period: 1 * time.Minute, Limit: 5}
	authLimiter := limiter.New(store, authRate, limiter.WithTrustForwardHeader(true))

	authGroup := r.Group("/")
	authGroup.Use(mgin.NewMiddleware(authLimiter))
	authGroup.Use(middleware.CSRF(isProd))
	{
		authGroup.POST("/register", ah.Register)
		authGroup.POST("/verify-otp", ah.VerifyOTP)
		authGroup.POST("/login", ah.Login)
		authGroup.POST("/logout", ah.Logout)
		authGroup.POST("/refresh", ah.Refresh)
		// OAuth endpoints are GET (browser redirect) and the callback is
		// excluded from CSRF in middleware/csrf.go since Google can't echo
		// our header. CSRF state is enforced via the oauth_state cookie.
		authGroup.GET("/auth/google/login", ah.GoogleLogin)
		authGroup.GET("/auth/google/callback", ah.GoogleCallback)
	}

	// ── Protected routes (JWT + Redis cache + CSRF on writes) ────────────────
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg, db))
	protected.Use(middleware.CSRF(isProd))

	guessRate := limiter.Rate{Period: 1 * time.Minute, Limit: 10}
	guessLimiter := limiter.New(store, guessRate, limiter.WithTrustForwardHeader(true))

	protected.Use(mgin.NewMiddleware(guessLimiter))
	protected.POST("/guess", ch.SubmitGuess)
	protected.GET("/challenge/attempts", ch.GetAttempts)
	protected.POST("/claim", ch.ClaimPrize)
	protected.GET("/challenge/my-status", ch.GetMyStatus)
}
