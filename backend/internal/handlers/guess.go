// femabras/backend/internal/handlers/guess.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	status := "incorrect"

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var challenge models.Challenge

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("release_date = ? AND is_active = ?", time.Now().UTC().Truncate(24*time.Hour), true).
			First(&challenge).Error; err != nil {
			return err
		}

		if len(req.Guess) != len(challenge.SecretCode) {
			return fmt.Errorf("guess must be exactly %d digits", len(challenge.SecretCode))
		}

		allowed := make(map[rune]bool)
		for _, r := range challenge.SecretCode {
			allowed[r] = true
		}
		for _, r := range req.Guess {
			if !allowed[r] {
				return fmt.Errorf("invalid digit")
			}
		}

		if req.Guess == challenge.SecretCode {
			status = "success"
			services.LockOnSuccess(context.Background(), userID, todayStr)

			var user models.User
			if err := tx.First(&user, userID).Error; err != nil {
				return fmt.Errorf("user not found for winner assignment")
			}

			winnerName := "Anonymous Hero"
			if user.Name != "" {
				winnerName = user.Name
			}

			return tx.Model(&challenge).Updates(map[string]interface{}{
				"is_active":   false,
				"winner_id":   userID,
				"winner_name": winnerName,
				"winner_pic":  user.Picture,
			}).Error
		}

		return nil
	})

	if err != nil {
		if err.Error() == "invalid digit" || err.Error() == "guess must be exactly..." {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active challenge available or already won"})
		}
		return
	}

	if status == "success" {
		remaining = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             status,
		"remaining_attempts": remaining,
	})
}

// ClaimPrize remains the same below...
func (h *ChallengeHandler) ClaimPrize(c *gin.Context) {
	var req models.ClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetString("user_id")

	var challenge models.Challenge
	if err := h.DB.Where("release_date = ?", time.Now().UTC().Truncate(24*time.Hour)).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	if challenge.WinnerID == nil || *challenge.WinnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the winner of this challenge"})
		return
	}

	if req.Method == "ATM" {
		if challenge.PrizeAmount < 1000 || challenge.PrizeAmount > 30000 || challenge.PrizeAmount%1000 != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This prize amount is not eligible for ATM withdrawal."})
			return
		}
	}

	payout := models.PayoutRequest{
		UserID:      userID,
		ChallengeID: challenge.ID,
		Amount:      challenge.PrizeAmount,
		Method:      req.Method,
		Destination: req.Destination,
		AccountName: req.AccountName,
		Status:      "pending",
	}

	if err := h.DB.Create(&payout).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Prize already claimed for this challenge"})
		return
	}

	var user models.User
	h.DB.First(&user, "id = ?", userID)
	username := "Anonymous Hero"
	if user.Name != "" {
		username = user.Name
	}

	go utils.SendAdminWinnerAlert(h.Cfg, username, challenge.PrizeAmount, req.Method)

	c.JSON(http.StatusOK, gin.H{"message": "Prize claim submitted successfully!"})
}
