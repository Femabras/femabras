// femabras/backend/internal/middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// userCacheTTL is set slightly under the 15-minute JWT access token TTL.
// A cached entry means the user existed and was verified when the token was
// last issued — safe to trust within this window.
const userCacheTTL = 14 * time.Minute

func Auth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		cookieToken, err := c.Cookie("access_token")
		if err == nil {
			tokenStr = cookieToken
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDAny := claims["user_id"]
		var userID string
		switch v := userIDAny.(type) {
		case string:
			userID = v
		case float64:
			userID = fmt.Sprintf("%.0f", v)
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type in token"})
			return
		}

		// ── Redis cache check ─────────────────────────────────────────────────
		// If the user's verified status is cached we skip the DB round-trip
		// entirely. Cache miss or Redis unavailable falls through to the DB.
		cacheKey := fmt.Sprintf("user:active:%s", userID)
		if cached, _ := services.GetCachedUserStatus(c.Request.Context(), cacheKey); cached {
			c.Set("user_id", userID)
			c.Next()
			return
		}

		// ── Database fallback ─────────────────────────────────────────────────
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User no longer exists"})
			return
		}

		if !user.IsVerified {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Account not verified"})
			return
		}

		// Cache the result so subsequent requests within the TTL skip the DB
		services.SetCachedUserStatus(c.Request.Context(), cacheKey, userCacheTTL)

		c.Set("user_id", userID)
		c.Next()
	}
}
