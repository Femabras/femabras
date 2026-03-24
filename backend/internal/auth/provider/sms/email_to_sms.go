// femabras/backend/internal/auth/provider/sms/email_to_sms.go
package sms

import (
	"context"
	"fmt"
	"github.com/Femabras/femabras/backend/internal/config"
	"net/smtp"
	"strings"
)

type EmailToSMSProvider struct {
	cfg *config.Config
}

func NewEmailToSMSProvider(cfg *config.Config) *EmailToSMSProvider {
	return &EmailToSMSProvider{cfg: cfg}
}

func (p *EmailToSMSProvider) Send(ctx context.Context, phone string, code string) error {
	phoneDigits := strings.TrimPrefix(phone, "+244")

	var gateway string
	// Unitel logic (91, 92, 93, 94, 95, 99)
	if strings.HasPrefix(phoneDigits, "91") || strings.HasPrefix(phoneDigits, "92") ||
		strings.HasPrefix(phoneDigits, "93") || strings.HasPrefix(phoneDigits, "94") {
		gateway = phoneDigits + "@sms.unitel.ao"
	} else if strings.HasPrefix(phoneDigits, "99") { // Movicel
		gateway = phoneDigits + "@sms.movicel.ao"
	} else {
		return fmt.Errorf("unsupported Angolan carrier for gateway")
	}

	auth := smtp.PlainAuth("", p.cfg.SMTPUsername, p.cfg.SMTPPassword, p.cfg.SMTPHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: OTP\r\n\r\nYour code: %s", gateway, code))

	return smtp.SendMail(p.cfg.SMTPHost+":"+p.cfg.SMTPPort, auth, p.cfg.FromEmail, []string{gateway}, msg)
}

func (p *EmailToSMSProvider) Verify(ctx context.Context, destination string, code string) (bool, error) {
	// Verification is handled by the database/service layer
	return true, nil
}
