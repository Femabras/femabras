// femabras/backend/internal/worker/processor.go
package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
)

// Processes the task when Asynq pulls it from Redis
func HandleVerificationEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailVerificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry) // Skip retry if JSON is corrupted
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️ RESEND_API_KEY is missing. Simulating email send for:", payload.Email)
		return nil
	}

	// 1. Build the Resend API Payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"from":    "Femabras Security <onboarding@resend.dev>", // Change to your verified domain later
		"to":      []string{payload.Email},
		"subject": "Your Femabras Secret Code",
		"html":    fmt.Sprintf("<strong>%s</strong> is your verification code. Welcome to the challenge!", payload.OTP),
	})
	if err != nil {
		return err
	}

	// 2. Execute Zero-Dependency HTTP Request
	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email via HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("email provider rejected request with status: %d", resp.StatusCode)
	}

	fmt.Printf("✅ Verification email securely dispatched to %s\n", payload.Email)
	return nil
}
