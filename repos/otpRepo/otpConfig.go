package otpRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"otp_service/internal/db"
)

// interface
// struct
// recivers funtions

type OTPConfigs interface {
	GetOTPConfig(ctx context.Context, clientID, otpType string) (*OtpConfigModel, error)
}

type otpConfig struct {
	db db.DB
}

func NewOTPConfigsRepo(db db.DB) OTPConfigs {
	return &otpConfig{db: db}
}

func (o *otpConfig) GetOTPConfig(ctx context.Context, clientID, otpType string) (*OtpConfigModel, error) {
	query := `SELECT 
            client_id, otp_type, reuse_otp_on_resend, 
            template_name, otp_length, expiry_seconds, max_verify_attempts, 
            resend_allowed, max_resend_attempts, cooldown_seconds, created_at, updated_at 
        FROM otp_configs
        WHERE client_id = $1 AND otp_type = $2`

	var OtpConfig OtpConfigModel
	err := o.db.GetContext(ctx, &OtpConfig, query, clientID, otpType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no OTP config found for client_id: %s and otp_type: %s", clientID, otpType)
		}
		return nil, err
	}

	return &OtpConfig, nil
}
