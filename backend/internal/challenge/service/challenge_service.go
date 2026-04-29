// femabras/backend/internal/challenge/service/challenge_service.go
package service

import (
	"context"
	"log"

	"github.com/Femabras/femabras/internal/challenge/repository"
	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/utils"
	"gorm.io/gorm"
)

type ChallengeService interface {
	ProcessGuess(ctx context.Context, userID string, guess string) (string, int, error)
}

type challengeService struct {
	repo repository.ChallengeRepository
}

func NewChallengeService(repo repository.ChallengeRepository) ChallengeService {
	return &challengeService{repo: repo}
}

func (s *challengeService) ProcessGuess(ctx context.Context, userID string, guess string) (string, int, error) {
	todayStr := utils.GetTodayDateString()
	todayTruncated := utils.GetTodayTruncated()

	remaining, err := services.DecrementAndSave(ctx, userID, todayStr)
	if err != nil {
		return "", 0, ErrServiceUnavailable
	}
	if remaining < 0 {
		return "", 0, ErrNoAttemptsLeft
	}

	status := "incorrect"

	err = s.repo.WithTransaction(func(tx *gorm.DB) error {
		challenge, err := s.repo.GetActiveChallengeForUpdate(tx, todayTruncated)
		if err != nil {
			return ErrChallengeNotFound
		}

		if len(guess) != len(challenge.SecretCode) {
			return ErrInvalidAtmAmount // reuse as generic validation — see handler
		}

		if guess == challenge.SecretCode {
			status = "success"
			services.LockOnSuccess(ctx, userID, todayStr)

			user, err := s.repo.GetUserByID(tx, userID)
			if err != nil {
				return ErrChallengeNotFound
			}

			winnerName := "Anonymous Hero"
			if user.Name != "" {
				winnerName = user.Name
			}

			if err := s.repo.UpdateChallengeAsWon(tx, challenge, userID, winnerName, user.Picture); err != nil {
				return err
			}

			if pubErr := services.PublishChallengeEvent(ctx, "solved"); pubErr != nil {
				log.Printf("warn: failed to publish challenge solved event: %v", pubErr)
			}
		}
		return nil
	})

	if err != nil {
		return "", remaining, err
	}

	if status == "success" {
		remaining = 0
	}

	return status, remaining, nil
}

// Keep models imported to avoid unused import error during migration
var _ = models.PayoutRequest{}
