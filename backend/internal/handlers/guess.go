// femabras/backend/internal/handlers/guess.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Femabras/femabras/backend/internal/models"
	"github.com/Femabras/femabras/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func (h *ChallengeHandler) SubmitGuess(c *gin.Context) {
	var req models.GuessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	todayStr := time.Now().UTC().Format("2006-01-02")

	// === REDIS + FALLBACK ===
	remaining, err := services.DecrementAndSave(c.Request.Context(), userID, todayStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}
	if remaining < 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "No attempts left today"})
		return
	}

	// === YOUR ORIGINAL VALIDATION ===
	var challenge models.Challenge // ← this fixes "undefined: challenge"
	if err := h.DB.Where("release_date = ? AND is_active = ?", time.Now().UTC().Truncate(24*time.Hour), true).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active challenge today"})
		return
	}

	if len(req.Guess) != len(challenge.SecretCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Guess must be exactly %d digits", len(challenge.SecretCode))})
		return
	}

	// digit validation...
	allowed := make(map[rune]bool)
	for _, r := range challenge.SecretCode {
		allowed[r] = true
	}
	for _, r := range req.Guess {
		if !allowed[r] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid digit"})
			return
		}
	}

	status := "incorrect"
	if req.Guess == challenge.SecretCode {
		status = "success"
		services.LockOnSuccess(context.Background(), userID, todayStr)
		h.DB.Model(&challenge).Update("is_active", false)
		remaining = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             status,
		"remaining_attempts": remaining,
	})
}
