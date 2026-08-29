package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to parse PostgreSQL config: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("❌ Failed to create PostgreSQL pool: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		log.Fatalf("❌ Failed to connect PostgreSQL: %v", err)
	}

	DB = db

	log.Println("✅ Connected to PostgreSQL!")
}
