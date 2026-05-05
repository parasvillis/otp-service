package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"otp_service/cmd"
	"otp_service/cmd/router"
	"otp_service/handlers"
	"otp_service/internal/db"
	"otp_service/internal/redisclient"
	"otp_service/repos"
	"otp_service/repos/otpRepo"
	"otp_service/services/contracts"
	"otp_service/services/delivery"
	otpservice "otp_service/services/otpService"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	dbConnStr := "postgres://postgres:postgres@db:5432/otp_service?sslmode=disable"
	ctx := context.Background()

	// dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_SSLMODE"),
	// )
	log.Printf("Connecting to Postgres with connection string: [%s]", dbConnStr)
	if dbConnStr == "" {
		log.Fatal("Database connection string is not set")
	}
	dbClient, err := db.NewDBClient(dbConnStr)
	if err != nil {
		log.Fatal("Error creating DB client: ", err)
	}

	redisAddr := "redis:6379"
	// redisAddr := os.Getenv("REDIS_ADDR")
	// log.Printf("Connecting to Redis at: [%s]", redisAddr)
	// if redisAddr == "" {
	// 	log.Fatal("REDIS_ADDR environment variable is not set")
	// }
	rClient := redisclient.NewRedisClient(redisAddr)
	if err := rClient.Ping(ctx); err != nil {
		log.Fatal("Redis ping err: ", err)
	}

	cmd.RunMigrations(dbConnStr)

	// pre build local or redis cache goes here if needed

	// infra
	rCache := repos.NewCache(rClient.RDB)

	// repos
	otpConfigRepo := otpRepo.NewOTPConfigsRepo(dbClient, logger)
	otpTemplateRepo := otpRepo.NewOtpTemplatesRepo(dbClient, logger)
	otpEventRepo := otpRepo.NewOtpEventRepo(dbClient, logger)

	// external services
	senderFactory := buildSenderFactory()

	// services
	otpService := otpservice.NewOptService(rCache, otpConfigRepo, otpTemplateRepo, otpEventRepo, senderFactory)

	// handlers
	container := router.HandlerContainer{
		OTP: handlers.NewOtpHandler(otpService),
	}
	r := router.InitRoutes(container)

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func buildSenderFactory() contracts.SenderFactory {
	return delivery.NewSenderFactory()
}

// fun c buildDBRepos(dbClient *db.Client, logger *zap.Logger) repos.Repo {
// 	otpConfigRepo := otpRepo.NewOTPConfigsRepo(dbClient, logger)
// 	otpTemplateRepo := otpRepo.NewOtpTemplatesRepo(dbClient, logger)
// 	otpEventRepo := otpRepo.NewOtpEventRepo(dbClient, logger)
// 	// cache := repos.GetCache(nil)
// 	return repos.NewRepoRegistry(otpConfigRepo, otpTemplateRepo, otpEventRepo, nil)
// }

// func registerJobs(serviceRegistry services.ServiceRegistry, logger *zap.Logger) {
// 	handler := jobHandler.NewJobsHandler(serviceRegistry, logger)
// 	handler.RegisterJobs(jobHandler.NewJob("test_job", 20, jobHandler.TestJob, true))
// 	handler.RunJobs()
// }
