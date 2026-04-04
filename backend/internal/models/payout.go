// backend/internal/models/payout.go
package models

import "gorm.io/gorm"

type PayoutRequest struct {
	gorm.Model
	UserID      string `gorm:"not null;index"`
	ChallengeID uint   `gorm:"not null;uniqueIndex"`
	Amount      int    `gorm:"not null"`
	Method      string `gorm:"not null"`
	Destination string `gorm:"not null"`
	AccountName string `gorm:"size:255"`
	Status      string `gorm:"default:'pending'"`
	AdminNotes  string
}

type ClaimRequest struct {
	Method      string `json:"method" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	AccountName string `json:"account_name"`
}
