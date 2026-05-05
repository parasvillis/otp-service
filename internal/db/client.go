package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	*sqlx.DB
}

func NewDBClient(dsn string) (*Client, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "pgx")

	return &Client{DB: dbx}, nil
}

func WaitForDB(url string) {
	for {
		db, err := sql.Open("postgres", url)
		if err == nil {
			if err = db.Ping(); err == nil {
				break
			}
		}
		log.Println("Waiting for DB...")
		time.Sleep(3 * time.Second)
	}
}
