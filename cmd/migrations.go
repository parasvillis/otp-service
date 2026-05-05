package cmd

import (
	"log"
	"otp_service/internal/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	db.WaitForDB(dbURL)
	m, err := migrate.New(
		"file://./migrations",
		// "file:///migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
