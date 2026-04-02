// femabras/backend/internal/auth/service/auth_service.go
package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Femabras/femabras/internal/auth/repository"
	"github.com/Femabras/femabras/internal/auth/types"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/models"
	"github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/worker"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req types.RegisterRequest) (uint, error)
	VerifyOTP(ctx context.Context, userID uint, code string) (string, error)
	HandleGoogleCallback(ctx context.Context, code string) (string, error)
	Login(ctx context.Context, req types.LoginRequest) (string, error)
	RunCleanup()
}

type authService struct {
	repo        repository.AuthRepository
	cfg         *config.Config
	asynqClient *asynq.Client
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config, asynqClient *asynq.Client) AuthService {
	return &authService{
		repo:        repo,
		cfg:         cfg,
		asynqClient: asynqClient,
	}
}

func (s *authService) Register(ctx context.Context, req types.RegisterRequest) (uint, error) {
	if _, err := s.repo.GetUserByIdentifier(req.Email); err == nil {
		return 0, errors.New("an account with this email already exists")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	otpCode := s.generateSecureOTP(6)
	hashedOTP, _ := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)

	pending := models.PendingUser{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashed),
		OTPCode:      string(hashedOTP),
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	if err := s.repo.CreatePendingUser(&pending); err != nil {
		return 0, errors.New("registration already in progress for this number")
	}

	if req.Email != "" {
		if s.asynqClient != nil {
			err := worker.EnqueueVerificationEmail(s.asynqClient, req.Email, otpCode)
			if err != nil {
				log.Printf("CRITICAL: Failed to queue email for %s: %v", req.Email, err)
			}
		}
	}

	return pending.ID, nil
}

func (s *authService) Login(ctx context.Context, req types.LoginRequest) (string, error) {
	user, err := s.repo.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", errors.New("please verify your phone number first")
	}

	todayStr := time.Now().UTC().Format("2006-01-02")
	_, _ = services.GetOrCreateAttempts(context.Background(), fmt.Sprintf("%d", user.ID), todayStr)

	return s.generateToken(user.ID)
}

func (s *authService) VerifyOTP(ctx context.Context, pendingID uint, code string) (string, error) {
	pending, err := s.repo.GetPendingUserByID(pendingID)
	if err != nil {
		return "", errors.New("invalid request")
	}
	if time.Now().After(pending.ExpiresAt) {
		return "", errors.New("OTP expired")
	}
	if pending.Attempts >= 5 {
		return "", errors.New("too many attempts")
	}

	err = bcrypt.CompareHashAndPassword([]byte(pending.OTPCode), []byte(code))
	if err != nil {
		_ = s.repo.IncrementOTPFail(pendingID)
		return "", errors.New("invalid OTP")
	}

	user := models.User{
		Name:         pending.Name,
		Email:        &pending.Email,
		PasswordHash: pending.PasswordHash,
		IsVerified:   true,
	}
	if pending.Email != "" {
		user.Email = &pending.Email
	}
	if err := s.repo.CreateUser(&user); err != nil {
		return "", err
	}

	todayStr := time.Now().UTC().Format("2006-01-02")
	_, _ = services.GetOrCreateAttempts(context.Background(), fmt.Sprintf("%d", user.ID), todayStr)

	_ = s.repo.DeletePendingUser(pending.ID)
	return s.generateToken(user.ID)
}

func GetGoogleAuthURL(state string) string {
	if googleOauthConfig == nil {
		return ""
	}
	return googleOauthConfig.AuthCodeURL(state)
}

func (s *authService) HandleGoogleCallback(ctx context.Context, code string) (string, error) {
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	client := googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", err
	}

	user, err := s.repo.GetUserByIdentifier(userInfo.Email)
	if err != nil {
		emailVal := userInfo.Email
		user = &models.User{
			GoogleID: &userInfo.ID,
			Email:    &emailVal,
			Name:     userInfo.Name,
		}
		if err := s.repo.CreateUser(user); err != nil {
			return "", err
		}
	}

	todayStr := time.Now().UTC().Format("2006-01-02")
	_, _ = services.GetOrCreateAttempts(context.Background(), fmt.Sprintf("%d", user.ID), todayStr)

	return s.generateToken(user.ID)
}

func (s *authService) generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) generateSecureOTP(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}

func (s *authService) RunCleanup() {
	if err := s.repo.CleanupExpiredPendingUsers(); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	if err := s.repo.ResetPendingSequence(); err != nil {
		log.Printf("Sequence reset error: %v", err)
	}
}
