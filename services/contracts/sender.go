package contracts

import (
	"context"
	"otp_service/dto"
)

type OtpSender interface {
	Send(ctx context.Context, req dto.SendOTPRequest) error
}

type SenderFactory interface {
	GetSender(channel, vendor string) (OtpSender, error)
}
