// femabras/backend/internal/handlers/contact.go
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Femabras/femabras/internal/config"
	"github.com/gin-gonic/gin"
)

type ContactRequest struct {
	Name    string `json:"name"    binding:"required,min=1,max=100"`
	Email   string `json:"email"   binding:"required,email"`
	Message string `json:"message" binding:"required,min=10,max=2000"`
}

// ContactHandler returns a Gin handler that accepts a contact form submission
// and forwards it to the admin via Resend. Falls back to SMTP if no Resend key.
func ContactHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ContactRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		// Sanitise inputs — strip leading/trailing whitespace
		req.Name = strings.TrimSpace(req.Name)
		req.Message = strings.TrimSpace(req.Message)

		if err := sendContactEmail(cfg, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message. Please try again."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message received. We'll be in touch soon."})
	}
}

func sendContactEmail(cfg *config.Config, req ContactRequest) error {
	if cfg.ResendAPIKey != "" {
		return sendViaResend(cfg, req)
	}
	// Fallback: log to stdout so the message is never silently lost
	fmt.Printf("CONTACT FORM (no Resend key) — From: %s <%s>\n%s\n", req.Name, req.Email, req.Message)
	return nil
}

func sendViaResend(cfg *config.Config, req ContactRequest) error {
	htmlBody := fmt.Sprintf(`
<h2>New Contact Form Submission</h2>
<p><strong>Name:</strong> %s</p>
<p><strong>Email:</strong> <a href="mailto:%s">%s</a></p>
<hr>
<p>%s</p>
`, req.Name, req.Email, req.Email,
		strings.ReplaceAll(req.Message, "\n", "<br>"))

	body, err := json.Marshal(map[string]interface{}{
		"from":     fmt.Sprintf("%s <noreply@auth.femabras.com>", cfg.AppName),
		"to":       []string{cfg.AdminEmail},
		"reply_to": req.Email,
		"subject":  fmt.Sprintf("[Contact] New message from %s", req.Name),
		"html":     htmlBody,
	})
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.ResendAPIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend rejected contact email with status %d", resp.StatusCode)
	}
	return nil
}
