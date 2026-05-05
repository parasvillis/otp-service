package helpers

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"otp_service/constants"
)

func GetOTPKey(userID, clientID, otpType string) string {
	return fmt.Sprintf("otp:%s:%s:%s", userID, clientID, otpType)
}

func GetIPAddress(ctx context.Context) *string {
	if ip, ok := ctx.Value(constants.ContextKeyIPAddress).(*string); ok {
		return ip
	}
	return nil
}

func GetUserAgent(ctx context.Context) *string {
	if ua, ok := ctx.Value(constants.ContextKeyUserAgent).(*string); ok {
		return ua
	}
	return nil
}

func GenerateSecureOTP(otpLength int) (string, error) {
	const digits = "0123456789"
	otpBytes := make([]byte, otpLength)
	for i := 0; i < otpLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otpBytes[i] = digits[num.Int64()]
	}
	return string(otpBytes), nil
}

func MapLuaError(code int64) error {
	switch code {
	case -1:
		return constants.ErrTooManyResends
	case -2:
		return constants.ErrOTPExpired
	case -3:
		return constants.ErrResendNotAllowed
	case -4:
		return constants.ErrStillInCooldown
	default:
		return constants.ErrUnknown
	}
}

func RenderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.
		New("otp").
		Option("missingkey=error").
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}
