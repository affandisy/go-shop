package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(databaseURI string, isProduction bool) error {
	var err error

	config := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	}

	if isProduction {
		config.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(postgres.Open(databaseURI), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
