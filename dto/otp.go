package dto

import (
	"fmt"
	"otp_service/utils"
	"strings"
)

type GenerateOTPRequest struct {
	UserID   string `json:"userID"`
	ClientID string `json:"clientID"`
	OtpType  string `json:"otpType"`
}

func GenerateOTPRequestValidator(req GenerateOTPRequest) error {
	req.UserID = strings.TrimSpace(req.UserID)
	req.ClientID = strings.TrimSpace(req.ClientID)
	req.OtpType = strings.TrimSpace(req.OtpType)

	if !utils.ValidateUUID(req.UserID) {
		return fmt.Errorf("invalid user id")
	}

	if !utils.ValidateUUID(req.ClientID) {
		return fmt.Errorf("invalid client id")
	}

	if req.OtpType == "" {
		return fmt.Errorf("otp type is required")
	}

	return nil
}

type GenerateOTPResponse struct {
	Messages []string
	Errors   []string
}

type ResendOtpRequest struct {
	GenerateOTPRequest
}

type ResendOTPResponse struct {
	GenerateOTPResponse
}

func ResendOTPRequestValidator(req ResendOtpRequest) error {
	return GenerateOTPRequestValidator(req.GenerateOTPRequest)
}

type VerifyOTPRequest struct {
	UserID   string `json:"userID"`
	ClientID string `json:"clientID"`
	OtpType  string `json:"otpType"`
	OTP      string `json:"otp"`
}

func VerifyOTPRequestValidator(req VerifyOTPRequest) error {
	req.UserID = strings.TrimSpace(req.UserID)
	req.ClientID = strings.TrimSpace(req.ClientID)
	req.OtpType = strings.TrimSpace(req.OtpType)
	req.OTP = strings.TrimSpace(req.OTP)

	if !utils.ValidateUUID(req.UserID) {
		return fmt.Errorf("invalid user id")
	}

	if !utils.ValidateUUID(req.ClientID) {
		return fmt.Errorf("invalid client id")
	}

	if req.OtpType == "" {
		return fmt.Errorf("otp type is required")
	}

	if req.OTP == "" {
		return fmt.Errorf("otp is required")
	}
	return nil
}

type OTPEventParams struct {
	ClientID      string
	UserID        string
	OtpType       string
	Identifier    string
	Channel       string
	EventType     string
	Status        string
	AttemptNumber int
	ResendNumber  int
	Vendor        string
	Error         *string
	Otp           *string
	SubmittedOtp  *string
	IPAddress     *string
	UserAgent     *string
}
