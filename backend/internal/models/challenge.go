// femabras/backend/internal/models/challenge.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	SecretCode  string    `gorm:"not null;size:32"`
	Difficulty  int       `gorm:"not null;default:1"` // 1=3 digits, 2=4, ..., up to 8
	ReleaseDate time.Time `gorm:"type:date;uniqueIndex;not null"`
	IsActive    bool      `gorm:"default:true"`
}
