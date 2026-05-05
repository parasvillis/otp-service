package dto

type SendOTPRequest struct {
	ClientID         string
	Identifier       string // phone or email
	Channel          string // sms/email/whatsapp
	OTP              string
	TemplateBody     string
	VendorTemplateID *string
	Vendor           string
	OTPType          string
	EmailSubject     *string
}
