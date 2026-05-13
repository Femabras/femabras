// femabras/backend/internal/middleware/security.go
package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets a strict Content Security Policy and other defensive
// headers on every response. Values are tuned for Femabras' specific needs:
//   - GAM rewarded ads load scripts from securepubads.g.doubleclick.net
//   - SSE connects to the same origin
//   - Google profile pictures from lh3.googleusercontent.com
func SecurityHeaders() gin.HandlerFunc {
	csp := "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline' 'unsafe-eval' " +
		"https://securepubads.g.doubleclick.net " +
		"https://pagead2.googlesyndication.com " +
		"https://tpc.googlesyndication.com; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: blob: " +
		"https://lh3.googleusercontent.com " +
		"https://*.googlesyndication.com " +
		"https://*.doubleclick.net; " +
		"font-src 'self' data:; " +
		"connect-src 'self' " +
		"https://securepubads.g.doubleclick.net " +
		"https://*.googlesyndication.com; " +
		"frame-src https://*.doubleclick.net https://*.googlesyndication.com; " +
		"object-src 'none'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"frame-ancestors 'none';"

	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", csp)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		// HSTS only meaningful in production over HTTPS — Railway terminates TLS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}
