// femabras/backend/internal/auth/provider/factory.go
package provider

import (
	"github.com/Femabras/femabras/backend/internal/auth/provider/email"
	"github.com/Femabras/femabras/backend/internal/auth/provider/sms"
	"github.com/Femabras/femabras/backend/internal/config"
)

type OTPFactory interface {
	GetProvider(method string) OTPProvider
}

type factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) OTPFactory {
	return &factory{cfg: cfg}
}

func (f *factory) GetProvider(method string) OTPProvider {
	switch method {
	case "twilio":
		return sms.NewTwilioProvider(f.cfg)
	case "email":
		return email.NewEmailProvider(f.cfg)
	case "gateway":
		return sms.NewEmailToSMSProvider(f.cfg)
	default:
		// Fallback to gateway if unspecified
		return sms.NewEmailToSMSProvider(f.cfg)
	}
}
