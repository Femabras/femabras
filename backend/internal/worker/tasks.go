// femabras/backend/internal/worker/tasks.go
package worker

import (
	"encoding/json"
	"fmt"

	"github.com/Femabras/femabras/internal/utils"
	"github.com/hibiken/asynq"
)

const (
	TypeSendVerificationEmail = "email:send_verification"
)

type EmailVerificationPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func EnqueueVerificationEmail(client *asynq.Client, email, otp string) error {
	rawPayload, err := json.Marshal(EmailVerificationPayload{Email: email, OTP: otp})
	if err != nil {
		return err
	}

	encryptedPayload, err := utils.EncryptPayload(rawPayload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendVerificationEmail, []byte(encryptedPayload))

	info, err := client.Enqueue(task, asynq.MaxRetry(3))
	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}

	fmt.Printf("Enqueued email task: id=%s queue=%s\n", info.ID, info.Queue)
	return nil
}
