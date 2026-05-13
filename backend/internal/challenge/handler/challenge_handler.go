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
	"gorm.io/gorm/clause"
)

type ChallengeHandler struct {
	svc service.ChallengeService
	db  *gorm.DB
	cfg *config.Config
}

func NewChallengeHandler(svc service.ChallengeService, db *gorm.DB, cfg *config.Config) *ChallengeHandler {
	return &ChallengeHandler{svc: svc, db: db, cfg: cfg}
}

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

	// Build the digit pool from the secret. The challenge struct only holds
	// the plaintext secret if this is a fresh creation; otherwise we'd need
	// to derive the pool differently. For now, we store digits in the same
	// transaction that creates the challenge — see services/challenge.go.
	//
	// Since the secret is hashed, we cannot recover digits from the hash.
	// Instead, we always include all digits 0-9 as the tray. This is a
	// gameplay simplification: every digit is theoretically available.
	// The challenge IS the position/sequence, not which digits are used.
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	sort.Strings(digits)

	c.JSON(http.StatusOK, gin.H{
		"status": "active",
		"slots":  challenge.Difficulty,
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
		case errors.Is(err, service.ErrInvalidGuessLength):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	}

	c.JSON(http.StatusOK, gin.H{
		"is_winner":     isWinner,
		"payout_status": payoutStatus,
	})
}

// ClaimPrize is now race-condition-safe via a SELECT FOR UPDATE on the
// challenge row inside a transaction. A concurrent double-click or replay
// will block on the row lock and find the payout already exists on the
// second attempt, returning a clear 200 with idempotency rather than 409.
func (h *ChallengeHandler) ClaimPrize(c *gin.Context) {
	var req models.ClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var alreadyExisted bool
	var prizeAmount int

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Lock the challenge row — concurrent claim attempts will serialise here
		var challenge models.Challenge
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("release_date = ?", utils.GetTodayTruncated()).
			First(&challenge).Error; err != nil {
			return service.ErrChallengeNotFound
		}

		if challenge.WinnerID == nil || *challenge.WinnerID != userID {
			return service.ErrNotWinner
		}

		// ATM amount validation (existing rule)
		if req.Method == "ATM" {
			if challenge.PrizeAmount < 1000 ||
				challenge.PrizeAmount > 30000 ||
				challenge.PrizeAmount%1000 != 0 {
				return service.ErrInvalidAtmAmount
			}
		}

		// Check for existing payout — idempotency.
		// If one exists for this user+challenge, return 200 OK without
		// creating a duplicate. This handles double-click and network retry.
		var existing models.PayoutRequest
		err := tx.Where("challenge_id = ? AND user_id = ?", challenge.ID, userID).
			First(&existing).Error
		if err == nil {
			alreadyExisted = true
			prizeAmount = challenge.PrizeAmount
			return nil
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("check existing payout: %w", err)
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

		if err := tx.Create(&payout).Error; err != nil {
			return service.ErrPrizeAlreadyClaimed
		}

		prizeAmount = challenge.PrizeAmount
		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, service.ErrChallengeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrNotWinner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidAtmAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrPrizeAlreadyClaimed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process claim"})
		}
		return
	}

	// Only fire the admin alert email on the first successful claim.
	// Repeated calls return 200 silently so the UI stays responsive.
	if !alreadyExisted {
		var user models.User
		h.db.First(&user, "id = ?", userID)
		username := "Anonymous Hero"
		if user.Name != "" {
			username = user.Name
		}
		go utils.SendAdminWinnerAlert(h.cfg, username, prizeAmount, req.Method)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Prize claim submitted successfully!",
		"already_claimed": alreadyExisted,
	})
}

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
