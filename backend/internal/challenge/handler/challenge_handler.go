// femabras/backend/internal/challenge/handler/challenge_handler.go
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Femabras/femabras/internal/challenge/service"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/models"
	appServices "github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChallengeHandler struct {
	svc service.ChallengeService
	db  *gorm.DB
	cfg *config.Config
}

func NewChallengeHandler(svc service.ChallengeService, db *gorm.DB, cfg *config.Config) *ChallengeHandler {
	return &ChallengeHandler{svc: svc, db: db, cfg: cfg}
}

// ── Game endpoints ────────────────────────────────────────────────────────────

func (h *ChallengeHandler) GetDailyChallenge(c *gin.Context) {
	challengePointer, err := appServices.CreateOrGetTodayChallenge(h.db)
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
			if err := h.db.Where("challenge_id = ?", challenge.ID).First(&payout).Error; err == nil {
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

	status, remaining, err := h.svc.ProcessGuess(c.Request.Context(), userID, req.Guess)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoAttemptsLeft):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrServiceUnavailable):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrChallengeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             status,
		"remaining_attempts": remaining,
	})
}

func (h *ChallengeHandler) GetAttempts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	todayStr := time.Now().UTC().Format("2006-01-02")
	attempts, err := appServices.GetOrCreateAttempts(c.Request.Context(), userID, todayStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attempts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"remaining_attempts": attempts})
}

func (h *ChallengeHandler) GetMyStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var challenge models.Challenge
	today := utils.GetTodayTruncated()
	if err := h.db.Where("release_date = ?", today).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	isWinner := false
	payoutStatus := "unclaimed"

	if challenge.WinnerID != nil && *challenge.WinnerID == userID {
		isWinner = true
		var payout models.PayoutRequest
		if err := h.db.Where("challenge_id = ? AND user_id = ?", challenge.ID, userID).First(&payout).Error; err == nil {
			payoutStatus = payout.Status
		}
		// No error log needed — "record not found" is a valid pre-claim state
	}

	c.JSON(http.StatusOK, gin.H{
		"is_winner":     isWinner,
		"payout_status": payoutStatus,
	})
}

func (h *ChallengeHandler) ClaimPrize(c *gin.Context) {
	var req models.ClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetString("user_id")

	var challenge models.Challenge
	if err := h.db.Where("release_date = ?", utils.GetTodayTruncated()).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	if challenge.WinnerID == nil || *challenge.WinnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": service.ErrNotWinner.Error()})
		return
	}

	if req.Method == "ATM" {
		if challenge.PrizeAmount < 1000 || challenge.PrizeAmount > 30000 || challenge.PrizeAmount%1000 != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrInvalidAtmAmount.Error()})
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

	if err := h.db.Create(&payout).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": service.ErrPrizeAlreadyClaimed.Error()})
		return
	}

	var user models.User
	h.db.First(&user, "id = ?", userID)
	username := "Anonymous Hero"
	if user.Name != "" {
		username = user.Name
	}

	go utils.SendAdminWinnerAlert(h.cfg, username, challenge.PrizeAmount, req.Method)

	c.JSON(http.StatusOK, gin.H{"message": "Prize claim submitted successfully!"})
}

// ── SSE ───────────────────────────────────────────────────────────────────────

// StreamStatus is a Server-Sent Events endpoint. It holds each browser
// connection open and pushes a single "solved" event the instant someone wins,
// replacing the previous 20-second polling loop entirely.
func (h *ChallengeHandler) StreamStatus(c *gin.Context) {
	rc := http.NewResponseController(c.Writer)
	if err := rc.SetWriteDeadline(time.Time{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx := c.Request.Context()
	pubsub := appServices.SubscribeChallengeEvents(ctx)
	defer pubsub.Close()

	ch := pubsub.Channel()

	fmt.Fprintf(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(c.Writer, "event: challenge\ndata: %s\n\n", msg.Payload)
			c.Writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(c.Writer, ": keepalive\n\n")
			c.Writer.Flush()
		case <-ctx.Done():
			return
		}
	}
}
