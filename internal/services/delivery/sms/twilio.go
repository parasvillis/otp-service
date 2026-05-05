package sms

import (
	"context"
	"fmt"
	"otp_service/dto"
	"time"
)

type TwilioSender struct {
	// client
}

func (m *TwilioSender) Send(ctx context.Context, req dto.SendOTPRequest) error {
	fmt.Println("Sending sms with twilio sender : ", req.OTP, req.TemplateBody)
	time.Sleep(2 * time.Second)
	return nil
}
