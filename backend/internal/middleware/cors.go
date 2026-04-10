// femabras/backend/internal/middleware/cors.go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOriginsStr string) gin.HandlerFunc {

	origins := strings.Split(allowedOriginsStr, ",")
	validOrigins := make([]string, 0)
	for _, o := range origins {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			validOrigins = append(validOrigins, trimmed)
		}
	}

	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")

		originToSet := ""
		if len(validOrigins) > 0 {
			originToSet = validOrigins[0]
		}

		for _, allowed := range validOrigins {
			if requestOrigin == allowed {
				originToSet = requestOrigin
				break
			}
		}

		if originToSet != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", originToSet)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
