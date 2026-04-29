// femabras/backend/internal/challenge/repository/challenge_repository.go
package repository

import (
	"github.com/Femabras/femabras/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ChallengeRepository interface {
	GetTodayChallenge(today time.Time) (*models.Challenge, error)
	GetActiveChallengeForUpdate(tx *gorm.DB, today time.Time) (*models.Challenge, error)
	UpdateChallengeAsWon(tx *gorm.DB, challenge *models.Challenge, userID string, userName string, userPic string) error
	CreatePayoutRequest(payout *models.PayoutRequest) error
	GetPayoutByChallengeAndUser(challengeID uint, userID string) (*models.PayoutRequest, error)
	GetUserByID(tx *gorm.DB, userID string) (*models.User, error)
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type challengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepository{db: db}
}

func (r *challengeRepository) GetTodayChallenge(today time.Time) (*models.Challenge, error) {
	var challenge models.Challenge
	err := r.db.Where("release_date = ?", today).First(&challenge).Error
	return &challenge, err
}

func (r *challengeRepository) GetActiveChallengeForUpdate(tx *gorm.DB, today time.Time) (*models.Challenge, error) {
	var challenge models.Challenge
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("release_date = ? AND is_active = ?", today, true).
		First(&challenge).Error
	return &challenge, err
}

func (r *challengeRepository) UpdateChallengeAsWon(tx *gorm.DB, challenge *models.Challenge, userID string, userName string, userPic string) error {
	return tx.Model(challenge).Updates(map[string]interface{}{
		"is_active":   false,
		"winner_id":   userID,
		"winner_name": userName,
		"winner_pic":  userPic,
	}).Error
}

func (r *challengeRepository) CreatePayoutRequest(payout *models.PayoutRequest) error {
	return r.db.Create(payout).Error
}

func (r *challengeRepository) GetPayoutByChallengeAndUser(challengeID uint, userID string) (*models.PayoutRequest, error) {
	var payout models.PayoutRequest
	err := r.db.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&payout).Error
	return &payout, err
}

func (r *challengeRepository) GetUserByID(tx *gorm.DB, userID string) (*models.User, error) {
	var user models.User
	err := tx.First(&user, userID).Error
	return &user, err
}

func (r *challengeRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
