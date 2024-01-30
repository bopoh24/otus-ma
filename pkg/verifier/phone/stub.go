package phone

import (
	"context"
	"errors"
	"strings"
)

var _ Verifier = (*StubPhoneVerify)(nil)

// StubPhoneVerify twilio sms verification instance
type StubPhoneVerify struct {
	phoneCodes map[string]struct{}
}

// NewStubPhoneVerify creates instance of Twilio verifier
func NewStubPhoneVerify() (*StubPhoneVerify, error) {
	return &StubPhoneVerify{
		phoneCodes: make(map[string]struct{}),
	}, nil
}

// Send sends verification code
func (t *StubPhoneVerify) Send(ctx context.Context, phone string) error {
	t.phoneCodes[phone] = struct{}{}
	return ctx.Err()
}

// Check checks verification code
func (t *StubPhoneVerify) Check(ctx context.Context, phone, code string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if _, ok := t.phoneCodes[phone]; !ok {
		return errors.New("phone not found")
	}
	// code should be a suffix of a phone in this stub
	if !strings.HasSuffix(phone, code) {
		return ErrIncorrectVerificationCode
	}
	return nil
}
