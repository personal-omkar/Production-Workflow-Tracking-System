package db

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once       sync.Once
	DBInstance *gorm.DB
)

// GetDB tries to connect to the PostgreSQL database with a retry mechanism
func GetDB() *gorm.DB {
	once.Do(func() {
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		// Validate environment variables
		if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" {
			slog.Warn("One or more required database environment variables are not set")
			return
		}

		// DSN (Data Source Name)
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

		const maxRetries = 5
		const retryInterval = 5 * time.Second // Adjust the retry interval as needed

		for i := 0; i < maxRetries; i++ {
			var err error
			DBInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				slog.Error("Failed to connect to database", "attempt", i+1, "error", err)
				time.Sleep(retryInterval)
				continue
			}
			slog.Info("Successfully connected to the database")
			return
		}
		slog.Error("Failed to connect to the database after multiple attempts", "attempts", maxRetries)
	})

	return DBInstance
}
