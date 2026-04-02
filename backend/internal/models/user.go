// femabras/backend/internal/models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        *string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	IsVerified   bool    `gorm:"default:false"`
	GoogleID     *string `gorm:"uniqueIndex"`
	Name         string
	Picture      string
}

type PendingUser struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	OTPCode      string    `gorm:"size:60"`
	Attempts     int       `gorm:"default:0"`
	ExpiresAt    time.Time `gorm:"index"`
}
