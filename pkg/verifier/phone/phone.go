package phone

import (
	"context"
)

type Verifier interface {
	Send(ctx context.Context, phone string) error
	Check(ctx context.Context, phone, code string) error
}
