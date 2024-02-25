package phone

import (
	"context"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var _ Verifier = (*TwilioPhoneVerify)(nil)

const statusApproved = "approved"

// TwilioPhoneVerify twilio sms verification instance
type TwilioPhoneVerify struct {
	client *twilio.RestClient
	srvSid string
}

// NewTwilioPhoneVerify creates instance of Twilio verifier
func NewTwilioPhoneVerify(friendlyName string) (*TwilioPhoneVerify, error) {
	client := twilio.NewRestClient()
	params := &verify.CreateServiceParams{}
	params.SetFriendlyName(friendlyName)

	resp, err := client.VerifyV2.CreateService(params)
	if err != nil {
		return nil, err
	}
	return &TwilioPhoneVerify{
		client: client,
		srvSid: *resp.Sid,
	}, nil
}

// Send sends verification code
func (t *TwilioPhoneVerify) Send(ctx context.Context, phone string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	params := &verify.CreateVerificationParams{}
	// must be in E.164 format
	params.SetTo(phone)
	params.SetChannel("sms")
	_, err := t.client.VerifyV2.CreateVerification(t.srvSid, params)
	return err
}

// Check checks verification code
func (t *TwilioPhoneVerify) Check(ctx context.Context, phone, code string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	params := &verify.CreateVerificationCheckParams{}
	// must be in E.164 format
	params.SetTo(phone)
	params.SetCode(code)

	resp, err := t.client.VerifyV2.CreateVerificationCheck(t.srvSid, params)
	if err != nil {
		return err
	}

	if *resp.Status != statusApproved {
		return ErrIncorrectVerificationCode
	}
	return nil
}
