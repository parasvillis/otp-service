package constants

import "errors"

var (
	ErrOTPExpired      = errors.New("otp expired or not found")
	ErrInvalidOTP      = errors.New("invalid otp")
	ErrTooManyAttempts = errors.New("too many attempts")
	ErrOTPMissing      = errors.New("key found but otp missing")

	ErrTooManyResends   = errors.New("too many resends")
	ErrStillInCooldown  = errors.New("still in cooldown")
	ErrResendNotAllowed = errors.New("resend not allowed")
	ErrUnknown          = errors.New("unknown error")
)

const (
	LuaReuse    = "__REUSE__"
	LuaGenerate = "__GENERATE__"
)

type contextKey string

const (
	ContextKeyIPAddress contextKey = "ip_address"
	ContextKeyUserAgent contextKey = "user_agent"
)

const (
	EventGenerated = "GENERATED"
	EventSent      = "SENT"
	EventVerify    = "VERIFY"
	// EventFailed    = "FAILED"
	EventResend = "RESEND"

	StatusSuccess = "success"
	StatusFailure = "failure"
)
