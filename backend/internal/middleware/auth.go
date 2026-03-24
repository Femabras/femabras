// femabras/backend/internal/middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Femabras/femabras/backend/internal/config"
	"github.com/Femabras/femabras/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Auth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenStr string

		// 1. Try cookie first (PRIMARY - secure)
		cookieToken, err := c.Cookie("auth_token")
		if err == nil {
			tokenStr = cookieToken
		} else {
			// 2. Fallback to Authorization header (for mobile / API clients)
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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

		// --- ACTIVE DATABASE CHECK ---
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			// This catches if the user was deleted from the DB
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User no longer exists"})
			return
		}

		// Ensures users are verified before letting them guess
		if !user.PhoneVerified {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Phone not verified"})
			return
		}
		// ----------------------------------

		c.Set("user_id", userID)
		c.Next()
	}
}
