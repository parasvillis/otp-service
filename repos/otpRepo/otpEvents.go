package otpRepo

import (
	"context"
	"otp_service/internal/db"

	"go.uber.org/zap"
)

type OtpEvent interface {
	RegisterOtpEvent(ctx context.Context, event OtpEventModel) error
}

type otpEvent struct {
	db     db.DB
	logger *zap.Logger
}

func NewOtpEventRepo(db db.DB, logger *zap.Logger) OtpEvent {
	return &otpEvent{db: db, logger: logger}
}

func (o *otpEvent) RegisterOtpEvent(ctx context.Context, event OtpEventModel) error {
	query := `INSERT INTO otp_events (id, client_id, user_id, otp_type, channel, identifier, error, event_type, status, attempt_number, resend_number, vendor, otp, submitted_otp, ip_address)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := o.db.ExecContext(ctx, query,
		event.ID,
		event.ClientID,
		event.UserID,
		event.OtpType,
		event.Channel,
		event.Identifier,
		event.Error,
		event.EventType,
		event.Status,
		event.AttemptNumber,
		event.ResendNumber,
		event.Vendor,
		event.Otp,
		event.SubmittedOtp,
		event.IPAddress,
	)
	return err
}
