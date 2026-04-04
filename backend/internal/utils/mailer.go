// femabras/backend/internal/utils/mailer.go
package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/Femabras/femabras/internal/config"
)

func SendAdminWinnerAlert(cfg *config.Config, username string, amount int, method string) {

	if cfg.SMTPHost == "" || cfg.AdminEmail == "" {
		fmt.Printf("🚨 ADMIN ALERT (Console Only): %s won %d AOA via %s\n", username, amount, method)
		return
	}

	// Format the email headers correctly to avoid spam filters
	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: 🚨 FEMABRÁS PAYOUT REQUEST\n", cfg.FromEmail, cfg.AdminEmail)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`Alert!

A user has just successfully cracked today's Secret Code and requested a payout.

Winner Alias: %s
Prize Amount: %d AOA
Requested Method: %s

Please log in to your secure Railway database to view the destination details (IBAN/Phone) and process the payout.

- %s Automated System
`, username, amount, method, cfg.AppName)

	message := []byte(header + mime + body)

	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	log.Printf("📧 Dialing SMTP Server %s...", addr)

	err := smtp.SendMail(addr, auth, cfg.SMTPUsername, []string{cfg.AdminEmail}, message)
	if err != nil {
		fmt.Printf("Failed to send admin email alert: %v\n", err)
	}
	log.Printf("✅ Winner Alert Email successfully sent to %s", cfg.AdminEmail)
}
