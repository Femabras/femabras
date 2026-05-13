// femabras/backend/internal/models/challenge.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Challenge stores the daily code-cracking puzzle.
//
// IMPORTANT: SecretCode is no longer persisted to the database.
// SecretCodeHash holds a bcrypt hash of the secret. The plaintext
// secret only exists in memory at challenge generation time and is
// discarded immediately after hashing. Guesses are validated by
// bcrypt comparison, never by string equality.
type Challenge struct {
	gorm.Model
	// SecretCode is a TRANSIENT field — never written to the database.
	// `gorm:"-"` excludes it from all GORM operations.
	SecretCode string `gorm:"-" json:"-"`

	// SecretCodeHash is the bcrypt hash of the secret. Length 60 = bcrypt output.
	SecretCodeHash string `gorm:"size:60;not null" json:"-"`

	Difficulty  int       `gorm:"not null"`
	PrizeAmount int       `gorm:"not null"`
	ReleaseDate time.Time `gorm:"uniqueIndex;not null"`
	IsActive    bool      `gorm:"default:true"`

	// Winner fields populated when challenge is solved
	WinnerID   *string
	WinnerName string
	WinnerPic  string
}

// GuessRequest is the body of POST /guess
type GuessRequest struct {
	Guess string `json:"guess" binding:"required,len=1|min=3,max=20,numeric"`
}

// ClaimRequest is the body of POST /claim — fully validated to prevent
// injection, oversized fields, and malformed payment destinations.
type ClaimRequest struct {
	Method      string `json:"method"       binding:"required,oneof=Bank ATM Multicaixa"`
	Destination string `json:"destination"  binding:"required,min=5,max=50"`
	AccountName string `json:"account_name" binding:"required,min=2,max=100"`
}
