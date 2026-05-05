package email

import (
	"context"
	"fmt"
	"otp_service/dto"
	"time"
)

type SendGridSender struct {
	// client
}

func (m *SendGridSender) Send(ctx context.Context, req dto.SendOTPRequest) error {
	fmt.Println("Sending email with sendGrid sender : ", req.OTP, req.TemplateBody)
	time.Sleep(2 * time.Second)
	return nil
}
