package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"otp_service/constants"
	"otp_service/dto"
	otpservice "otp_service/services/otpService"
	"otp_service/utils"
)

type OtpHandler struct {
	otpService otpservice.OTPService
}

func NewOtpHandler(otp otpservice.OTPService) *OtpHandler {
	return &OtpHandler{otpService: otp}
}

func (o *OtpHandler) GenerateOTPHandler(w http.ResponseWriter, r *http.Request) {
	// request extraction and validation
	var req dto.GenerateOTPRequest
	ipAddress := utils.GetClientIP(r)
	userAgent := r.UserAgent()
	ctx := r.Context()

	ctx = context.WithValue(ctx, "ip_address", ipAddress)
	ctx = context.WithValue(ctx, "user_agent", userAgent)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := dto.GenerateOTPRequestValidator(req); err != nil {
		http.Error(w, "Request Validation Failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// call service layer
	res, err := o.otpService.Generate(ctx, req)
	if err != nil {
		http.Error(w, "error in otp generation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	httpStatus := http.StatusOK
	if len(res.Errors) > 0 {
		httpStatus = http.StatusInternalServerError
	}
	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{"res": res})
}

func (o *OtpHandler) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOTPRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := dto.VerifyOTPRequestValidator(req); err != nil {
		http.Error(w, "request validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := o.otpService.Verify(r.Context(), req)
	if err != nil {
		switch err {
		case constants.ErrInvalidOTP:
			http.Error(w, "invalid otp", http.StatusUnauthorized)
			return
		case constants.ErrTooManyAttempts:
			http.Error(w, "too many attempts", http.StatusForbidden)
			return
		case constants.ErrOTPExpired:
			http.Error(w, "otp expired", http.StatusGone)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if !valid {
		http.Error(w, "invalid otp", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"OTP verified successfully"}`))
}

func (o *OtpHandler) ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ResendOtpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := dto.ResendOTPRequestValidator(req); err != nil {
		http.Error(w, "request validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	res, err := o.otpService.Resend(r.Context(), req)
	if err != nil {
		http.Error(w, "error in resending otp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	httpStatus := http.StatusOK
	if len(res.Errors) > 0 {
		httpStatus = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{"res": res})

}
