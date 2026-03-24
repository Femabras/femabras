// femabras/backend/internal/auth/provider/sms/twilio.go
package sms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Femabras/femabras/backend/internal/config"
)

type TwilioProvider struct {
	cfg *config.Config
}

func NewTwilioProvider(cfg *config.Config) *TwilioProvider {
	return &TwilioProvider{cfg: cfg}
}

func (p *TwilioProvider) Send(ctx context.Context, destination string, code string) error {
	// Twilio API Endpoint
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", p.cfg.TwilioSID)

	// Message body
	msgBody := fmt.Sprintf("Your %s verification code is: %s", p.cfg.AppName, code)

	// Form Data
	data := url.Values{}
	data.Set("To", destination)
	data.Set("From", p.cfg.TwilioPhone)
	data.Set("Body", msgBody)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Basic Auth
	req.SetBasicAuth(p.cfg.TwilioSID, p.cfg.TwilioAuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("twilio API error: status %d", resp.StatusCode)
	}

	return nil
}

func (p *TwilioProvider) Verify(ctx context.Context, destination string, code string) (bool, error) {
	// Verification is handled by our internal database logic
	return true, nil
}
