// femabras/backend/internal/handlers/challenge.go
package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/Femabras/femabras/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChallengeHandler struct {
	DB *gorm.DB
}

func (h *ChallengeHandler) GetDailyChallenge(c *gin.Context) {

	challengePointer, err := services.CreateOrGetTodayChallenge(h.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load daily challenge"})
		return
	}

	challenge := *challengePointer

	if !challenge.IsActive {
		c.JSON(http.StatusOK, gin.H{
			"status":  "solved",
			"message": "Today's challenge has already been conquered!",
			"winner": gin.H{
				"name":    challenge.WinnerName,
				"picture": challenge.WinnerPic,
			},
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
