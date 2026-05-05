package delivery

import (
	"fmt"
	"otp_service/internal/services/contracts"
	"otp_service/internal/services/delivery/email"
	"otp_service/internal/services/delivery/sms"
)

// type OtpSender interface {
// 	Send(ctx context.Context, req dto.SendOTPRequest) error
// }

type SenderFactory struct {
	msg91Sender    contracts.OtpSender
	twilioSender   contracts.OtpSender
	sendGridSender contracts.OtpSender
}

func NewSenderFactory() contracts.SenderFactory {
	return &SenderFactory{
		msg91Sender:    &sms.Msg91Sender{},
		twilioSender:   &sms.TwilioSender{},
		sendGridSender: &email.SendGridSender{},
	}
}

func (f *SenderFactory) GetSender(channel, vendor string) (contracts.OtpSender, error) {
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
