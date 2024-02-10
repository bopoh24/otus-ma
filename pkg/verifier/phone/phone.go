package phone

import (
	"context"
)

//go:generate mockgen -source phone.go -destination mocks/verifier.go -package mock Verifier
type Verifier interface {
	Send(ctx context.Context, phone string) error
	Check(ctx context.Context, phone, code string) error
}
