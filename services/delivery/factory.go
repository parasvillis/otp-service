package delivery

import (
	"context"
	"fmt"
	"otp_service/dto"
	"otp_service/services/delivery/email"
	"otp_service/services/delivery/sms"
)

type OtpSender interface {
	Send(ctx context.Context, req dto.SendOTPRequest) error
}

type SenderFactory struct {
	msg91Sender    OtpSender
	twilioSender   OtpSender
	sendGridSender OtpSender
}

func NewSenderFactory() *SenderFactory {
	return &SenderFactory{
		msg91Sender:    &sms.Msg91Sender{},
		twilioSender:   &sms.TwilioSender{},
		sendGridSender: &email.SendGridSender{},
	}
}

func (f *SenderFactory) GetSender(channel, vendor string) (OtpSender, error) {
	switch channel {
	case "sms":
		switch vendor {
		case "msg91":
			return f.msg91Sender, nil
		case "twilio":
			return f.twilioSender, nil
		}
	case "email":
		switch vendor {
		case "sendgrid":
			return f.sendGridSender, nil
		}
	}
	return nil, fmt.Errorf("unsupported channel/vendor combination: %s/%s", channel, vendor)
}
