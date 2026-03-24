// femabras/backend/internal/models/otp.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Code      string    `gorm:"not null;size:6"`
	ExpiresAt time.Time `gorm:"not null"`
}
