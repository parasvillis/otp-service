package otpRepo

import (
	"time"
)

type OtpConfigModel struct {
	ClientID          string    `json:"client_id" db:"client_id"`
	OtpType           string    `json:"otp_type" db:"otp_type"`
	ReuseOtpOnResend  bool      `json:"reuse_otp_on_resend" db:"reuse_otp_on_resend"`
	TemplateName      string    `json:"template_name" db:"template_name"`
	OtpLength         int       `json:"otp_length" db:"otp_length"`
	ExpirySeconds     int       `json:"expiry_seconds" db:"expiry_seconds"`
	MaxVerifyAttempts *int      `json:"max_verify_attempts" db:"max_verify_attempts"`
	ResendAllowed     bool      `json:"resend_allowed" db:"resend_allowed"`
	MaxResendAttempts *int      `json:"max_resend_attempts" db:"max_resend_attempts"`
	CooldownSeconds   *int      `json:"cooldown_seconds" db:"cooldown_seconds"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Vendor       string `json:"vendor" db:"vendor"`               // from otp_templates
	// Channel      string `json:"channel" db:"channel"`             // from otp_templates
	// TemplateBody string `json:"template_body" db:"template_body"` // from otp_templates
}

type OtpTemplateModel struct {
	ID               string    `json:"id" db:"id"`
	ClientID         string    `json:"client_id" db:"client_id"`
	OtpType          string    `json:"otp_type" db:"otp_type"`
	Channel          string    `json:"channel" db:"channel"`
	TemplateName     string    `json:"template_name" db:"template_name"`
	TemplateBody     string    `json:"template_body" db:"template_body"`
	VendorTemplateID *string   `json:"vendor_template_id" db:"vendor_template_id"`
	Vendor           string    `json:"vendor" db:"vendor"`
	EmailSubject     *string   `json:"email_subject" db:"email_subject"`
	ISActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type OtpEventModel struct {
	ID            string    `db:"id"`
	ClientID      string    `db:"client_id"`
	UserID        string    `db:"user_id"`
	OtpType       string    `db:"otp_type"`
	Identifier    string    `db:"identifier"`
	Error         *string   `db:"error"`
	Channel       string    `db:"channel"`
	EventType     string    `db:"event_type"`
	Status        string    `db:"status"`
	AttemptNumber int       `db:"attempt_number"`
	ResendNumber  int       `db:"resend_number"`
	Vendor        string    `db:"vendor"`
	Otp           *string   `db:"otp"`
	SubmittedOtp  *string   `db:"submitted_otp"`
	IPAddress     *string   `db:"ip_address"`
	UserAgent     *string   `db:"user_agent"`
	CreatedAt     time.Time `db:"created_at"`
}
