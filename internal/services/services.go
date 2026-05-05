package services

// import (
// 	"otp_service/services/delivery"
// 	otpservice "otp_service/services/otpService"
// )

// type ServiceRegistry interface {
// 	GetOTPService() otpservice.OTPService
// 	GetDeliveryFactory() delivery.SenderFactory
// }

// type serviceRegistry struct {
// 	otpService      otpservice.OTPService
// 	deliveryFactory delivery.SenderFactory
// }

// func (s *serviceRegistry) GetOTPService() otpservice.OTPService {
// 	return s.otpService
// }

// func (s *serviceRegistry) GetDeliveryFactory() delivery.SenderFactory {
// 	return s.deliveryFactory
// }

// func NewServiceRegistry(otpService otpservice.OTPService, deliveryFactory delivery.SenderFactory) ServiceRegistry {
// 	return &serviceRegistry{
// 		otpService:      otpService,
// 		deliveryFactory: deliveryFactory,
// 	}
// }
