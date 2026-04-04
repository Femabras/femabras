// femabras/backend/internal/handlers/challenge.go
package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChallengeHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func (h *ChallengeHandler) GetDailyChallenge(c *gin.Context) {

	challengePointer, err := services.CreateOrGetTodayChallenge(h.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load daily challenge"})
		return
	}

	challenge := *challengePointer
	userID := c.GetString("user_id")

	if !challenge.IsActive {
		isWinner := false
		payoutStatus := "unclaimed"

		if userID != "" && challenge.WinnerID != nil && *challenge.WinnerID == userID {
			isWinner = true
			var payout models.PayoutRequest

			if err := h.DB.Where("challenge_id = ?", challenge.ID).First(&payout).Error; err == nil {
				payoutStatus = payout.Status
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "solved",
			"message": "Today's challenge has already been conquered!",
			"prize":   challenge.PrizeAmount,
			"winner": gin.H{
				"name":    challenge.WinnerName,
				"picture": challenge.WinnerPic,
			},
			"is_winner":     isWinner,
			"payout_status": payoutStatus,
		})
		return
	}

	digitSet := make(map[rune]bool)
	for _, r := range challenge.SecretCode {
		digitSet[r] = true
	}

	digits := make([]string, 0, len(digitSet))
	for r := range digitSet {
		digits = append(digits, string(r))
	}

	sort.Strings(digits)

	c.JSON(http.StatusOK, gin.H{
		"status": "active",
		"slots":  len(challenge.SecretCode),
		"date":   challenge.ReleaseDate.Format("2006-01-02"),
		"digits": digits,
		"prize":  challenge.PrizeAmount,
	})
}

func (h *ChallengeHandler) GetAttempts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	todayStr := time.Now().UTC().Format("2006-01-02")

	attempts, err := services.GetOrCreateAttempts(c.Request.Context(), userID, todayStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attempts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"remaining_attempts": attempts,
	})
}

func (h *ChallengeHandler) GetMyStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var challenge models.Challenge
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := h.DB.Where("release_date = ?", today).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	isWinner := false
	payoutStatus := "unclaimed"

	if challenge.WinnerID != nil && *challenge.WinnerID == userID {
		isWinner = true
		var payout models.PayoutRequest
		err := h.DB.Where("challenge_id = ? AND user_id = ?", challenge.ID, userID).First(&payout).Error
		if err == nil {
			payoutStatus = payout.Status
		} else {
			fmt.Printf("🚨 GET STATUS ERROR: Found winner, but failed to find PayoutRequest in DB: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"is_winner":     isWinner,
		"payout_status": payoutStatus,
	})
}
