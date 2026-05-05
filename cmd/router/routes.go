package router

import (
	"otp_service/internal/handlers"

	"github.com/go-chi/chi"
)

type HandlerContainer struct {
	OTP *handlers.OtpHandler
}

func InitRoutes(hc HandlerContainer) chi.Router {
	r := chi.NewRouter()
	r.Mount("/otp", otpRouterV1(hc.OTP))
	return r
}
