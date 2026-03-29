// femabras/backend/internal/handlers/guess.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"

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

	remaining, err := services.DecrementAndSave(c.Request.Context(), userID, todayStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}
	if remaining < 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "No attempts left today"})
		return
	}

	var challenge models.Challenge

	if err := h.DB.Where("release_date = ? AND is_active = ?", time.Now().UTC().Truncate(24*time.Hour), true).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active challenge today"})
		return
	}

	if len(req.Guess) != len(challenge.SecretCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Guess must be exactly %d digits", len(challenge.SecretCode))})
		return
	}

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

		// Fetch the winning user's details
		var user models.User
		h.DB.First(&user, userID)

		// Respect privacy: Use their Name, fallback to "Anonymous Hero"
		winnerName := "Anonymous Hero"

		if user.Name != "" {
			winnerName = user.Name
		}

		// Update the challenge with the winner's info
		h.DB.Model(&challenge).Updates(map[string]interface{}{
			"is_active":   false,
			"winner_id":   userID,
			"winner_name": winnerName,
			"winner_pic":  user.Picture,
		})

		remaining = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             status,
		"remaining_attempts": remaining,
	})
}
