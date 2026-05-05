package router

import (
	"otp_service/internal/handlers"

	"github.com/go-chi/chi"
)

func otpRouterV1(h *handlers.OtpHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/generate", h.GenerateOTPHandler)
	r.Post("/verify", h.VerifyOTPHandler)
	r.Post("/resend", h.ResendOTPHandler)
	return r
}
