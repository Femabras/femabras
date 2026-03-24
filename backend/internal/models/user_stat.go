// femabras/backend/internal/models/user_stat.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStat struct {
	gorm.Model
	UserID            string    `gorm:"not null;index:idx_user_date,unique" json:"user_id"`
	RemainingAttempts int       `gorm:"default:5;not null" json:"remaining_attempts"`
	AdsWatchedToday   int       `gorm:"default:0;not null" json:"ads_watched_today"`
	HintsUnlocked     bool      `gorm:"default:false" json:"hints_unlocked"`
	LastUpdate        time.Time `gorm:"type:date;not null;index:idx_user_date,unique" json:"last_update"`
}

// Request for guess submission (validation moved to handler)
type GuessRequest struct {
	Guess string `json:"guess" binding:"required"`
}
