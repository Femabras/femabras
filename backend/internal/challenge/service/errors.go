// femabras/backend/internal/challenge/service/errors.go
package service

import "errors"

// Sentinel errors for the challenge domain.
// Handlers use errors.Is() against these — never .Error() string matching.
var (
	ErrNoAttemptsLeft      = errors.New("no attempts left today")
	ErrServiceUnavailable  = errors.New("service temporarily unavailable")
	ErrChallengeNotFound   = errors.New("no active challenge available")
	ErrAlreadyWon          = errors.New("challenge already won")
	ErrNotWinner           = errors.New("you are not the winner of this challenge")
	ErrPrizeAlreadyClaimed = errors.New("prize already claimed for this challenge")
	ErrInvalidAtmAmount    = errors.New("this prize amount is not eligible for ATM withdrawal")
	ErrInvalidGuessLength  = errors.New("guess length does not match challenge")
)
