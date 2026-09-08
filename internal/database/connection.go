package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is not set")
	}

	db, err := gorm.Open(
		gorm_postgres.Open(databaseDSN),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
