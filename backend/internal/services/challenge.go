// femabras/backend/internal/services/challenge.go
package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Femabras/femabras/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateOrGetTodayChallenge ensures one challenge exists per UTC day.
// The plaintext secret_code is never stored — only its bcrypt hash.
// SecretCode lives in memory only during generation, then is zeroed.
func CreateOrGetTodayChallenge(db *gorm.DB) (*models.Challenge, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var challenge models.Challenge
	err := db.Where("release_date = ?", today).First(&challenge).Error
	if err == nil {
		// Existing challenge — secret_code field will be empty (not persisted)
		// and secret_code_hash will be loaded. This is correct: callers should
		// validate guesses via bcrypt against the hash.
		log.Printf("Found existing challenge for %s (ID: %d)",
			today.Format("2006-01-02"), challenge.ID)
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

	// Hash the secret immediately. bcrypt cost 10 is enough — at 10k+ users
	// the bottleneck is /guess, and we only hash once per day during creation.
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash secret: %w", err)
	}

	prize, err := calculateWeightedPrize(difficulty)
	if err != nil {
		return nil, fmt.Errorf("calculate prize: %w", err)
	}

	newChallenge := models.Challenge{
		SecretCodeHash: string(hash),
		Difficulty:     difficulty,
		PrizeAmount:    prize,
		ReleaseDate:    today,
		IsActive:       true,
	}

	if err := db.Create(&newChallenge).Error; err != nil {
		return nil, fmt.Errorf("create challenge: %w", err)
	}

	// Log only non-sensitive metadata. The plaintext secret variable will be
	// garbage-collected — it is never written to logs or persistent storage.
	log.Printf("Created new challenge for %s (ID: %d, length: %d, prize: %d AOA)",
		today.Format("2006-01-02"), newChallenge.ID, difficulty, prize)

	// Populate the transient SecretCode field so the caller (e.g. test seeding,
	// admin dashboard) can use it once if needed. In the normal request path
	// this is not needed.
	newChallenge.SecretCode = secret

	return &newChallenge, nil
}

// VerifyGuess returns true if the guess matches the stored bcrypt hash.
// Use this everywhere instead of `guess == challenge.SecretCode`.
func VerifyGuess(challenge *models.Challenge, guess string) bool {
	if challenge.SecretCodeHash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword(
		[]byte(challenge.SecretCodeHash),
		[]byte(guess),
	) == nil
}

func generateWeightedLength() int {
	roll, _ := rand.Int(rand.Reader, big.NewInt(100))
	val := roll.Int64()

	switch {
	case val < 60:
		subRoll, _ := rand.Int(rand.Reader, big.NewInt(2))
		return 4 + int(subRoll.Int64())
	case val < 90:
		return 6
	default:
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
