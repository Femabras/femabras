// femabras/backend/internal/middleware/csrf.go
package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRF protection via the double-submit cookie pattern.
//
// On any GET request we issue a non-HttpOnly "csrf_token" cookie containing
// a random per-session token. On state-changing requests (POST/PUT/PATCH/DELETE)
// we require the same token in the X-CSRF-Token header. The browser's same-origin
// policy prevents cross-site JavaScript from reading the cookie value, while
// any genuine same-origin request can read the cookie and echo it in the header.
//
// This is layered on top of SameSite=Lax cookies for defence in depth.

const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
	csrfTokenSize  = 32 // 256 bits of entropy
)

// CSRF returns a middleware that enforces the double-submit pattern.
// Apply it to any route group that mutates state via cookie-based auth.
func CSRF(isProd bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Issue a CSRF token cookie on safe methods if one doesn't exist.
		// The frontend reads this cookie via JavaScript and echoes it in the
		// X-CSRF-Token header on subsequent state-changing requests.
		if isSafeMethod(c.Request.Method) {
			ensureCSRFCookie(c, isProd)
			c.Next()
			return
		}

		// State-changing request — verify the token.
		// Skip the OAuth callback (Google can't send our header) and the
		// ad webhook (HMAC-signed by ad network, not browser-driven).
		if c.Request.URL.Path == "/auth/google/callback" ||
			c.Request.URL.Path == "/webhooks/ad-reward" {
			c.Next()
			return
		}

		cookieToken, err := c.Cookie(csrfCookieName)
		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token missing. Refresh the page and try again.",
			})
			return
		}

		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token header missing",
			})
			return
		}

		// Constant-time comparison prevents timing oracle on the token.
		if subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch",
			})
			return
		}

		c.Next()
	}
}

// IssueCSRFToken is a lightweight handler that ensures the cookie exists.
// Useful as a dedicated GET /csrf endpoint the frontend can call on app boot.
func IssueCSRFToken(isProd bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ensureCSRFCookie(c, isProd)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

func ensureCSRFCookie(c *gin.Context, isProd bool) {
	if existing, err := c.Cookie(csrfCookieName); err == nil && existing != "" {
		return
	}

	b := make([]byte, csrfTokenSize)
	if _, err := rand.Read(b); err != nil {
		// Cryptographic failure — extremely unlikely; fail closed
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSRF token",
		})
		return
	}
	token := hex.EncodeToString(b)

	domain := ""
	sameSite := http.SameSiteLaxMode
	if isProd {
		domain = ".femabras.com"
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 days, matches refresh token lifetime
		HttpOnly: false,            // MUST be readable by JavaScript
		Secure:   isProd,
		Domain:   domain,
		SameSite: sameSite,
	})
}
