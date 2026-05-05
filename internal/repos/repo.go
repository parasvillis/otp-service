package repos

// import (
// 	"otp_service/internal/db"
// 	"otp_service/repos/otpRepo"

// 	"github.com/redis/go-redis/v9"
// 	"go.uber.org/zap"
// )

// type Repo interface {
// 	GetOTPConfigRepo(db db.DB, logger *zap.Logger) otpRepo.OTPConfigs
// 	GetOTPTemplateRepo(db db.DB, logger *zap.Logger) otpRepo.OtpTemplate
// 	GetOTPEventRepo(db db.DB, logger *zap.Logger) otpRepo.OtpEvent
// 	GetCache(c *redis.Client) Cache
// }

// type repoRegistory struct {
// 	otpConfigRepo   otpRepo.OTPConfigs
// 	otpTemplateRepo otpRepo.OtpTemplate
// 	otpEventRepo    otpRepo.OtpEvent
// 	cache           Cache
// }

// func (s *repoRegistory) GetOTPConfigRepo(db db.DB, logger *zap.Logger) otpRepo.OTPConfigs {
// 	return s.otpConfigRepo
// }

// func (s *repoRegistory) GetOTPTemplateRepo(db db.DB, logger *zap.Logger) otpRepo.OtpTemplate {
// 	return s.otpTemplateRepo
// }

// func (s *repoRegistory) GetOTPEventRepo(db db.DB, logger *zap.Logger) otpRepo.OtpEvent {
// 	return s.otpEventRepo
// }

// func (s *repoRegistory) GetCache(c *redis.Client) Cache {
// 	return s.cache
// }

// func NewRepoRegistry(otpConfigRepo otpRepo.OTPConfigs, otpTemplateRepo otpRepo.OtpTemplate, otpEventRepo otpRepo.OtpEvent, cache Cache) Repo {
// 	return &repoRegistory{
// 		otpConfigRepo:   otpConfigRepo,
// 		otpTemplateRepo: otpTemplateRepo,
// 		otpEventRepo:    otpEventRepo,
// 		cache:           cache,
// 	}
// }
