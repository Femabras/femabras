// femabras/backend/internal/auth/provider/email/pure_email.go
package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/Femabras/femabras/backend/internal/config"
)

type EmailProvider struct {
	cfg *config.Config
}

func NewEmailProvider(cfg *config.Config) *EmailProvider {
	return &EmailProvider{cfg: cfg}
}

func (e *EmailProvider) Send(ctx context.Context, to string, code string) error {
	auth := smtp.PlainAuth("", e.cfg.SMTPUsername, e.cfg.SMTPPassword, e.cfg.SMTPHost)

	// Format the email properly with headers
	subject := fmt.Sprintf("Subject: %s Verification Code\r\n", e.cfg.AppName)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><h2>Your verification code is: <strong>%s</strong></h2></body></html>", code)

	msg := []byte(subject + mime + body)

	return smtp.SendMail(e.cfg.SMTPHost+":"+e.cfg.SMTPPort, auth, e.cfg.FromEmail, []string{to}, msg)
}

func (e *EmailProvider) Verify(ctx context.Context, destination string, code string) (bool, error) {
	return true, nil // verification done in service layer
}
