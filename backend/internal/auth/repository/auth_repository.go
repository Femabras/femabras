// femabras/backend/internal/auth/repository/auth_repository.go
package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Femabras/femabras/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByIdentifier(identifier string) (*models.User, error)
	UpdateUserStatus(id uint, verified bool) error

	CreatePendingUser(pending *models.PendingUser) error
	GetPendingUserByID(id uint) (*models.PendingUser, error)
	DeletePendingUser(id uint) error
	IncrementOTPFail(id uint) error

	CleanupExpiredPendingUsers() error
	ResetPendingSequence() error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreatePendingUser(pending *models.PendingUser) error {

	r.db.Unscoped().Where("email = ?", pending.Email).Delete(&models.PendingUser{})

	return r.db.Create(pending).Error
}

func (r *authRepository) GetPendingUserByID(id uint) (*models.PendingUser, error) {
	var pending models.PendingUser

	err := r.db.Where("id = ?", id).First(&pending).Error
	if err != nil {
		return nil, err
	}

	return &pending, nil
}

func (r *authRepository) IncrementOTPFail(id uint) error {
	return r.db.Model(&models.PendingUser{}).Where("id = ?", id).UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}

func (r *authRepository) DeletePendingUser(id uint) error {
	return r.db.Unscoped().Delete(&models.PendingUser{}, id).Error
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *authRepository) UpdateUserStatus(id uint, verified bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_verified", verified).Error
}

func (r *authRepository) CreateOTP(otp *models.OTP) error {
	return r.db.Create(otp).Error
}

func (r *authRepository) GetValidOTP(userID uint, code string) (*models.OTP, error) {
	var otp models.OTP

	err := r.db.Where("user_id = ? AND code = ? AND expires_at > NOW()", userID, code).First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, err
	}
	return &otp, nil
}

func (r *authRepository) DeleteOTP(id uint) error {
	return r.db.Unscoped().Delete(&models.OTP{}, id).Error
}

func (r *authRepository) CleanupExpiredPendingUsers() error {
	var acquired bool

	err := r.db.Raw(
		"SELECT pg_try_advisory_lock(hashtext('cleanup_pending_users'))",
	).Scan(&acquired).Error

	if err != nil || !acquired {
		return nil
	}

	defer func() {
		if err := r.db.Exec(
			"SELECT pg_advisory_unlock(hashtext('cleanup_pending_users'))",
		).Error; err != nil {
			log.Printf("failed to release advisory lock: %v", err)
		}
	}()

	result := r.db.Unscoped().
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&models.PendingUser{})

	if result.Error != nil {
		return fmt.Errorf("janitor cleanup failed: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		log.Printf("Janitor cleaned %d expired users", result.RowsAffected)
	}

	return nil
}

func (r *authRepository) ResetPendingSequence() error {
	var count int64
	r.db.Model(&models.PendingUser{}).Count(&count)
	if count == 0 {

		return r.db.Exec("ALTER SEQUENCE pending_users_id_seq RESTART WITH 1").Error
	}
	return nil
}

func (r *authRepository) GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", identifier).First(&user).Error
	return &user, err
}
