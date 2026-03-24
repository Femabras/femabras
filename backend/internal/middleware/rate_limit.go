// femabras/backend/internal/middleware/rate_limit.go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5, // 5 requests per minute per IP (adjust as needed)
	}

	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	return mgin.NewMiddleware(rateLimiter)
}
