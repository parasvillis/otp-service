package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"otp_service/cmd"
	"otp_service/handlers"
	"otp_service/internal/db"
	"otp_service/internal/redisclient"
	"otp_service/repos"
	"otp_service/repos/otpRepo"
	"otp_service/services/delivery"
	otpservice "otp_service/services/otpService"

	"github.com/go-chi/chi"
)

func main() {
	// dbConnStr := "postgres://postgres:postgres@localhost:5432/otp_service?sslmode=disable"
	ctx := context.Background()
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	log.Printf("Connecting to Postgres with connection string: [%s]", dbConnStr)
	if dbConnStr == "" {
		log.Fatal("Database connection string is not set")
	}
	dbClient, err := db.NewDBClient(dbConnStr)
	if err != nil {
		log.Fatal("Error creating DB client: ", err)
	}

	// redisAddr := "localhost:6379"
	redisAddr := os.Getenv("REDIS_ADDR")
	log.Printf("Connecting to Redis at: [%s]", redisAddr)
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}
	rClient := redisclient.NewRedisClient(redisAddr)
	if err := rClient.Ping(ctx); err != nil {
		log.Fatal("Redis ping err: ", err)
	}

	cmd.RunMigrations(dbConnStr)

	r := chi.NewRouter()

	rCache := repos.NewCache(rClient.RDB)
	otpConfigRepo := otpRepo.NewOTPConfigsRepo(dbClient)
	otpTemplateRepo := otpRepo.NewOtpTemplatesRepo(dbClient)
	otpEventRepo := otpRepo.NewOtpEventRepo(dbClient)
	senderFactory := delivery.NewSenderFactory()
	otpService := otpservice.NewOptService(rCache, otpConfigRepo, otpTemplateRepo, otpEventRepo, senderFactory)
	otpHandler := handlers.NewOtpHandler(otpService)

	r.Route("/otp", func(r chi.Router) {
		r.Post("/generate", otpHandler.GenerateOTPHandler)
		r.Post("/verify", otpHandler.VerifyOTPHandler)
		r.Post("/resend", otpHandler.ResendOTPHandler)
	})

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
