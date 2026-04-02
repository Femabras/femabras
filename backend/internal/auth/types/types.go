// femabras/backend/internal/auth/types/types.go
package types

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=16"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	OTP    string `json:"otp" binding:"required,len=6"`
}
