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

	difficulty := generateWeightedLength()

	secret, difficulty, err := generateRandomSecret(difficulty, difficulty)
	if err != nil {
		return nil, fmt.Errorf("generate secret: %w", err)
	}

	prize, err := calculateWeightedPrize(difficulty)
	if err != nil {
		return nil, fmt.Errorf("calculate prize: %w", err)
	}

	newChallenge := models.Challenge{
		SecretCode:  secret,
		Difficulty:  difficulty,
		PrizeAmount: prize,
		ReleaseDate: today,
		IsActive:    true,
	}

	if err := db.Create(&newChallenge).Error; err != nil {
		return nil, fmt.Errorf("create challenge: %w", err)
	}

	log.Printf("Created new challenge for %s: %s (length %d, prize %d AOA)", today.Format("2006-01-02"), secret, difficulty, prize)

	return &newChallenge, nil
}

func generateWeightedLength() int {
	// Roll a number between 0 and 99
	roll, _ := rand.Int(rand.Reader, big.NewInt(100))
	val := roll.Int64()

	switch {
	case val < 60:
		// 60% chance: Length 4 or 5 (The standard game)
		subRoll, _ := rand.Int(rand.Reader, big.NewInt(2))
		return 4 + int(subRoll.Int64())
	case val < 90:
		// 30% chance: Length 6 (The hard game)
		return 6
	default:
		// 10% chance: Length 7 (The extreme jackpot game)
		return 7
	}
}

func calculateWeightedPrize(length int) (int, error) {
	min, max := 0, 0

	switch {
	case length <= 3:
		min, max = 50, 200
	case length >= 4 && length <= 5:
		min, max = 200, 700
	case length >= 6 && length <= 8:
		min, max = 700, 2000
	case length >= 9 && length <= 16:
		min, max = 2000, 10000
	default:
		min, max = 10000, 30000
	}

	diff := int64(max - min + 1)

	r1, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}
	r2, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}

	offset := r1.Int64()
	if r2.Int64() < offset {
		offset = r2.Int64()
	}

	return min + int(offset), nil
}

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
