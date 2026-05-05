package sms

import (
	"context"
	"fmt"
	"otp_service/dto"
	"time"
)

type Msg91Sender struct {
	// client
}

func (m *Msg91Sender) Send(ctx context.Context, req dto.SendOTPRequest) error {
	fmt.Println("Sending sms with msg91 sender : ", req.OTP, req.TemplateBody)
	time.Sleep(2 * time.Second)
	return nil
}
