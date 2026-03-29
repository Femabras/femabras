// femabras/backend/internal/services/challenge.go
package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Femabras/femabras/internal/models"
	"gorm.io/gorm"
)

// CreateOrGetTodayChallenge ensures a challenge exists for today
func CreateOrGetTodayChallenge(db *gorm.DB) (*models.Challenge, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var challenge models.Challenge
	err := db.Where("release_date = ?", today).First(&challenge).Error
	if err == nil {
		log.Printf("Found existing challenge for %s: %s", today.Format("2006-01-02"), challenge.SecretCode)
		return &challenge, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("check challenge: %w", err)
	}

	// Generate new random challenge
	secret, difficulty, err := generateRandomSecret(3, 8)
	if err != nil {
		return nil, fmt.Errorf("generate secret: %w", err)
	}

	newChallenge := models.Challenge{
		SecretCode:  secret,
		Difficulty:  difficulty,
		ReleaseDate: today,
		IsActive:    true,
	}

	if err := db.Create(&newChallenge).Error; err != nil {
		return nil, fmt.Errorf("create challenge: %w", err)
	}

	log.Printf("Created new challenge for %s: %s (length %d)", today.Format("2006-01-02"), secret, difficulty)
	return &newChallenge, nil
}

// generateRandomSecret creates a random digit string between min and max length
func generateRandomSecret(minLen, maxLen int) (string, int, error) {
	lenRange := int64(maxLen - minLen + 1)
	lengthBig, err := rand.Int(rand.Reader, big.NewInt(lenRange))
	if err != nil {
		return "", 0, err
	}
	length := minLen + int(lengthBig.Int64())

	var builder []byte
	for i := 0; i < length; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", 0, err
		}
		builder = append(builder, byte('0'+digit.Int64()))
	}

	return string(builder), length, nil
}
