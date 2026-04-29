// femabras/backend/internal/handlers/ads.go
package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/services"
	"github.com/gin-gonic/gin"
)

type AdRewardPayload struct {
	UserID    string `form:"user_id" binding:"required"`
	Reward    int    `form:"reward_amount" binding:"required"`
	TransID   string `form:"transaction_id" binding:"required"`
	Timestamp int64  `form:"timestamp" binding:"required"`
	Hash      string `form:"hash" binding:"required"`
}

// AdRewardWebhook returns a Gin handler that validates ad network S2S callbacks.
// The network secret is read from cfg — never hardcoded.
func AdRewardWebhook(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload AdRewardPayload
		if err := c.ShouldBindQuery(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		now := time.Now().Unix()
		if now-payload.Timestamp > 300 || payload.Timestamp-now > 60 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Payload expired or invalid timestamp"})
			return
		}

		expectedData := payload.TransID + payload.UserID + strconv.FormatInt(payload.Timestamp, 10)
		h := hmac.New(sha256.New, []byte(cfg.AdNetworkSecret))
		h.Write([]byte(expectedData))
		expectedHash := hex.EncodeToString(h.Sum(nil))

		if !hmac.Equal([]byte(payload.Hash), []byte(expectedHash)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		todayStr := time.Now().UTC().Format("2006-01-02")

		if err := services.IncrementAttemptAndAdsWatched(c.Request.Context(), payload.UserID, todayStr); err != nil {
			if err.Error() == "daily ad limit reached" {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Daily ad limit reached for user"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant reward"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "reward granted"})
	}
}
