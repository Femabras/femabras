// femabras/backend/internal/auth/provider/provider.go
package provider

import "context"

type OTPProvider interface {
	Send(ctx context.Context, destination string, code string) error
	Verify(ctx context.Context, destination string, code string) (bool, error)
}
