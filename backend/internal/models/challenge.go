// femabras/backend/internal/models/challenge.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	SecretCode  string    `gorm:"not null;size:32"`
	Difficulty  int       `gorm:"not null;default:1"`
	ReleaseDate time.Time `gorm:"type:date;uniqueIndex;not null"`
	IsActive    bool      `gorm:"default:true"`
	WinnerID    *string   `gorm:"index"`
	WinnerName  *string   `gorm:"size:255"`
	WinnerPic   *string   `gorm:"size:255"`
}
