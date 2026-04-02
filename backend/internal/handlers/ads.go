// femabras/backend/internal/handlers/ads.go
package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/Femabras/femabras/internal/services"
	"github.com/gin-gonic/gin"
)

const AdNetworkSecret = "your_super_secret_s2s_key"

type AdRewardPayload struct {
	UserID  string `form:"user_id" binding:"required"`
	Reward  int    `form:"reward_amount" binding:"required"`
	TransID string `form:"transaction_id" binding:"required"`
	Hash    string `form:"hash" binding:"required"`
}

func AdRewardWebhook(c *gin.Context) {
	var payload AdRewardPayload
	if err := c.ShouldBindQuery(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	expectedData := payload.TransID + payload.UserID
	h := hmac.New(sha256.New, []byte(AdNetworkSecret))
	h.Write([]byte(expectedData))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(payload.Hash), []byte(expectedHash)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature format"})
		return
	}

	todayStr := time.Now().UTC().Format("2006-01-02")

	err := services.IncrementAttemptAndAdsWatched(c.Request.Context(), payload.UserID, todayStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grant reward"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "reward granted"})
}
