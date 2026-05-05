package otpservice

import (
	"context"
	"fmt"
	"log"
	"otp_service/constants"
	"otp_service/dto"
	"otp_service/internal/repos"
	"otp_service/internal/repos/otpRepo"
	"otp_service/internal/services/contracts"
	"otp_service/internal/services/helpers"
	"otp_service/utils"

	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPService interface {
	Generate(ctx context.Context, req dto.GenerateOTPRequest) (dto.GenerateOTPResponse, error)
	Verify(ctx context.Context, req dto.VerifyOTPRequest) (bool, error)
	Resend(ctx context.Context, req dto.ResendOtpRequest) (dto.ResendOTPResponse, error)
}

type otpService struct {
	cache           repos.Cache
	otpConfigRepo   otpRepo.OTPConfigs
	otpTemplateRepo otpRepo.OtpTemplate
	otpEventRepo    otpRepo.OtpEvent
	senderFactory   contracts.SenderFactory
}

func NewOptService(c repos.Cache, otpConfigRepo otpRepo.OTPConfigs, otpTemplateRepo otpRepo.OtpTemplate, otpEventRepo otpRepo.OtpEvent, senderFactory contracts.SenderFactory) OTPService {
	return &otpService{cache: c, otpConfigRepo: otpConfigRepo, otpTemplateRepo: otpTemplateRepo, otpEventRepo: otpEventRepo, senderFactory: senderFactory}
}

func (o *otpService) Generate(ctx context.Context, req dto.GenerateOTPRequest) (dto.GenerateOTPResponse, error) {
	userID := req.UserID
	clientID := req.ClientID
	otpType := req.OtpType
	ipAddress := helpers.GetIPAddress(ctx)
	userAgent := helpers.GetUserAgent(ctx)
	var response dto.GenerateOTPResponse

	otpConfig, err := o.otpConfigRepo.GetOTPConfig(ctx, clientID, otpType)

	if err != nil {
		return response, fmt.Errorf("failed to fetch OTP config: %w", err)
	}
	if otpConfig == nil {
		return response, fmt.Errorf("no OTP config found for client_id: %s and otp_type: %s", clientID, otpType)
	}

	otp, err := helpers.GenerateSecureOTP(otpConfig.OtpLength)
	if err != nil {
		o.logOTPEvent(ctx, dto.OTPEventParams{
			ClientID:   clientID,
			UserID:     userID,
			OtpType:    otpType,
			Identifier: "1234567890",
			Channel:    "",
			EventType:  constants.EventGenerated,
			Status:     constants.StatusFailure,
			Error:      utils.StringPtr(err.Error()),
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
		})
		return response, fmt.Errorf("failed to generate OTP: %w", err)
	}

	key := helpers.GetOTPKey(userID, clientID, otpType)

	fmt.Println("key: ", key)

	// Store OTP and reset attempts in a hash
	if err := o.cache.HSet(ctx, key, map[string]interface{}{
		"otp":                 otp,
		"verify_attempts":     0,
		"max_verify_attempts": utils.IntPtrValue(otpConfig.MaxVerifyAttempts, -1),
		"resend_allowed":      otpConfig.ResendAllowed,
		"resend_attempts":     0,
		"max_resend_attempts": utils.IntPtrValue(otpConfig.MaxResendAttempts, -1),
		"reuse_otp_on_resend": otpConfig.ReuseOtpOnResend,
		"cooldown_seconds":    utils.IntPtrValue(otpConfig.CooldownSeconds, 10),
		"otp_type":            otpConfig.OtpType,
		"last_sent":           time.Now().Unix(),
	}).Err(); err != nil {
		return response, fmt.Errorf("error setting OTP in Redis: %w", err)
	}

	// Set TTL for the whole hash
	if err := o.cache.Expire(ctx, key, time.Duration(otpConfig.ExpirySeconds)*time.Second).Err(); err != nil {
		log.Printf("failed to set TTL for OTP key: %v", err)
	}
	o.logOTPEvent(ctx, dto.OTPEventParams{
		ClientID:      clientID,
		UserID:        userID,
		OtpType:       otpType,
		Identifier:    "1234567890",
		Channel:       "",
		EventType:     constants.EventGenerated,
		Status:        constants.StatusSuccess,
		AttemptNumber: 0,
		ResendNumber:  0,
		Vendor:        "",
		Otp:           &otp,
		SubmittedOtp:  nil,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	})
	o.sendOTP(ctx, otp, &response, otpConfig, userID)

	return response, nil
}

func (o *otpService) Verify(ctx context.Context, req dto.VerifyOTPRequest) (bool, error) {
	ipAddress := helpers.GetIPAddress(ctx)
	userAgent := helpers.GetUserAgent(ctx)
	key := helpers.GetOTPKey(req.UserID, req.ClientID, req.OtpType)

	result, err := o.cache.RunLuaInt(ctx, verifyOtpLuaScript, []string{key}, req.OTP)
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}

	status := constants.StatusFailure
	eventType := constants.EventVerify

	switch result {
	case 1:
		status = constants.StatusSuccess
	case 0:
		err = constants.ErrInvalidOTP
	case -1:
		err = constants.ErrTooManyAttempts
	case -2:
		err = constants.ErrOTPExpired
	case -3:
		err = constants.ErrOTPMissing
	default:
		err = fmt.Errorf("unknown OTP verification result: %d", result)
	}

	verifyAttemptNumber, _ := o.cache.HGet(ctx, key, "verify_attempts").Result()
	verifyAttemptNumberInt, _ := strconv.Atoi(verifyAttemptNumber)
	resendAttemptNumber, _ := o.cache.HGet(ctx, key, "resend_attempts").Result()
	resendAttemptNumberInt, _ := strconv.Atoi(resendAttemptNumber)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	o.logOTPEvent(ctx, dto.OTPEventParams{
		ClientID:      req.ClientID,
		UserID:        req.UserID,
		OtpType:       req.OtpType,
		Identifier:    "1234567890",
		EventType:     eventType,
		Status:        status,
		AttemptNumber: verifyAttemptNumberInt,
		ResendNumber:  resendAttemptNumberInt,
		SubmittedOtp:  &req.OTP,
		Error:         utils.StringPtr(errStr),
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	})

	if result == 1 {
		return true, nil
	}
	return false, err
}

func (o *otpService) Resend(ctx context.Context, req dto.ResendOtpRequest) (dto.ResendOTPResponse, error) {
	clientID := req.ClientID
	otpType := req.OtpType
	userID := req.UserID
	ipAddress := helpers.GetIPAddress(ctx)
	userAgent := helpers.GetUserAgent(ctx)

	var response dto.ResendOTPResponse
	otpConfig, err := o.otpConfigRepo.GetOTPConfig(ctx, clientID, otpType)
	if err != nil {
		return response, fmt.Errorf("failed to fetch OTP config: %w", err)
	}
	if otpConfig == nil {
		return response, fmt.Errorf("no OTP config found for client_id: %s and otp_type: %s", clientID, otpType)
	}

	otpKey := helpers.GetOTPKey(userID, clientID, otpType)
	res, err := o.cache.RunLuaAny(ctx, resendOtpLuaScript, []string{otpKey}, time.Now().Unix())
	if err != nil {
		return response, fmt.Errorf("failed to run resend OTP Lua: %w", err)
	}
	verifyAttemptNumber, _ := o.cache.HGet(ctx, otpKey, "verify_attempts").Result()
	verifyAttemptNumberInt, _ := strconv.Atoi(verifyAttemptNumber)
	resendAttemptNumber, _ := o.cache.HGet(ctx, otpKey, "resend_attempts").Result()
	resendAttemptNumberInt, _ := strconv.Atoi(resendAttemptNumber)
	otp, err := o.handleLuaResponse(ctx, res, otpKey, otpConfig)

	if err != nil {
		o.logOTPEvent(ctx, dto.OTPEventParams{
			ClientID:      req.ClientID,
			UserID:        req.UserID,
			OtpType:       req.OtpType,
			Error:         utils.StringPtr(err.Error()),
			Identifier:    "1234567890",
			EventType:     constants.EventResend,
			Status:        constants.StatusFailure,
			AttemptNumber: verifyAttemptNumberInt,
			ResendNumber:  resendAttemptNumberInt,
			Otp:           &otp,
			IPAddress:     ipAddress,
			UserAgent:     userAgent,
		})
		return response, err
	}
	o.sendOTP(ctx, otp, &response.GenerateOTPResponse, otpConfig, userID)

	o.logOTPEvent(ctx, dto.OTPEventParams{
		ClientID:      req.ClientID,
		UserID:        req.UserID,
		OtpType:       req.OtpType,
		Identifier:    "1234567890",
		EventType:     constants.EventResend,
		Status:        constants.StatusSuccess,
		AttemptNumber: verifyAttemptNumberInt,
		ResendNumber:  resendAttemptNumberInt + 1,
		Otp:           &otp,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	})
	return response, nil
}

func (o *otpService) handleLuaResponse(ctx context.Context, res interface{}, otpKey string, cfg *otpRepo.OtpConfigModel) (string, error) {
	var otp string
	var err error
	switch v := res.(type) {
	case int64:
		return "", helpers.MapLuaError(v)

	case string:
		switch v {
		case constants.LuaReuse:
			otp, err = o.cache.HGet(ctx, otpKey, "otp").Result()

		case constants.LuaGenerate:
			otp, err = helpers.GenerateSecureOTP(cfg.OtpLength)
		default:
			return "", fmt.Errorf("unexpected lua string response: %s", v)
		}

		if err != nil {
			return "", err
		}
		if err := o.postResendOTP(ctx, otpKey, otp, cfg.ExpirySeconds); err != nil {
			return "", err
		}

		return otp, nil

	default:
		return "", fmt.Errorf("unexpected Lua return type: %T", res)
	}
}

func (o *otpService) postResendOTP(ctx context.Context, otpKey string, otp string, ttlSeconds int) error {
	_, err := o.cache.RunLuaAny(ctx, resendPostLuaScript, []string{otpKey}, otp, ttlSeconds, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to update OTP atomically in Redis: %w", err)
	}

	return nil
}

func (o *otpService) sendOTP(ctx context.Context, otp string, response *dto.GenerateOTPResponse, otpConfig *otpRepo.OtpConfigModel, userID string) {

	ipAddress := helpers.GetIPAddress(ctx)
	userAgent := helpers.GetUserAgent(ctx)
	otpTemplates, err := o.otpTemplateRepo.GetOtpTemplates(ctx, otpConfig.ClientID, otpConfig.TemplateName)
	if err != nil {
		log.Printf("failed to fetch OTP template: %v", err)
		response.Errors = append(response.Errors, "failed to fetch OTP template: "+err.Error())
		return
	}
	if len(otpTemplates) == 0 {
		log.Printf("no OTP templates found for client %s and template %s",
			otpConfig.ClientID, otpConfig.TemplateName)
		response.Errors = append(response.Errors, "no OTP templates configured")
		return
	}

	// we can add caching for templates also.

	for _, tpl := range otpTemplates {
		log.Printf("Would send OTP using template: %s via channel: %s", tpl.TemplateName, tpl.Channel)
		sender, err := o.senderFactory.GetSender(tpl.Channel, tpl.Vendor)
		if err != nil {
			log.Printf("sender not found for channel %s: %v", tpl.Channel, err)
			response.Errors = append(response.Errors, err.Error())
			continue
		}
		renderedTemplate, err := helpers.RenderTemplate(tpl.TemplateBody, map[string]interface{}{
			"otp":            otp,
			"otp_type":       otpConfig.OtpType,
			"expiry_minutes": otpConfig.ExpirySeconds / 60,
		})
		if err != nil {
			log.Printf("failed to render template: %v", err)
			response.Errors = append(response.Errors, err.Error())
			continue
		}
		response.Messages = append(response.Messages, renderedTemplate)
		sendReq := dto.SendOTPRequest{
			ClientID:         otpConfig.ClientID,
			Identifier:       userID,
			Channel:          tpl.Channel,
			OTP:              otp,
			TemplateBody:     renderedTemplate,
			VendorTemplateID: tpl.VendorTemplateID,
			Vendor:           tpl.Vendor,
			OTPType:          otpConfig.OtpType,
			EmailSubject:     tpl.EmailSubject,
		}
		detachedCtx := context.WithoutCancel(ctx)
		go func(s contracts.OtpSender, sendReq dto.SendOTPRequest, tModel otpRepo.OtpTemplateModel) {
			sendCtx, cancel := context.WithTimeout(detachedCtx, 5*time.Second)
			defer cancel()
			status := constants.StatusSuccess
			var errStr *string

			if err := s.Send(sendCtx, sendReq); err != nil {
				log.Printf("failed to send OTP via %s: %v", tModel.Channel, err)
				status = constants.StatusFailure
				e := err.Error()
				errStr = &e
			}

			o.logOTPEvent(detachedCtx, dto.OTPEventParams{
				ClientID:      otpConfig.ClientID,
				UserID:        userID,
				OtpType:       otpConfig.OtpType,
				Identifier:    "1234567890",
				Channel:       tModel.Channel,
				EventType:     constants.EventSent,
				Status:        status,
				AttemptNumber: 0,
				ResendNumber:  0,
				Vendor:        tModel.Vendor,
				Error:         errStr,
				Otp:           &otp,
				SubmittedOtp:  nil,
				IPAddress:     ipAddress,
				UserAgent:     userAgent,
			})
		}(sender, sendReq, tpl)

	}
}

func (o *otpService) logOTPEvent(ctx context.Context, params dto.OTPEventParams) {
	event := otpRepo.OtpEventModel{
		ID:            utils.GenerateUUID(),
		ClientID:      params.ClientID,
		UserID:        params.UserID,
		OtpType:       params.OtpType,
		Identifier:    params.Identifier,
		Error:         params.Error,
		Channel:       params.Channel,
		EventType:     strings.ToUpper(params.EventType),
		Status:        strings.ToLower(params.Status),
		AttemptNumber: params.AttemptNumber,
		ResendNumber:  params.ResendNumber,
		Vendor:        params.Vendor,
		Otp:           params.Otp,
		SubmittedOtp:  params.SubmittedOtp,
		IPAddress:     params.IPAddress,
		UserAgent:     params.UserAgent,
	}
	detachedCtx := context.WithoutCancel(ctx)
	go func() {
		logCtx, cancel := context.WithTimeout(detachedCtx, 5*time.Second)
		defer cancel()

		if err := o.otpEventRepo.RegisterOtpEvent(logCtx, event); err != nil {
			log.Printf("failed to log OTP event: %v", err)
		}
	}()
}
