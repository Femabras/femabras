// femabras/backend/internal/database/database.go
package database

import (
	"log"
	"time"

	"github.com/Femabras/femabras/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Postgres defaults to max_connections=100.
	// With multiple backend instances at 10k+ scale, keep this conservative
	// and add PgBouncer in front of Postgres when scaling horizontally.
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	log.Println("Running AutoMigrate...")
	err = db.AutoMigrate(
		&models.PendingUser{},
		&models.Challenge{},
		&models.UserStat{},
		&models.User{},
		&models.OTP{},
		&models.PayoutRequest{},
		&models.RefreshToken{},
	)
	if err != nil {
		return nil, err
	}
	log.Println("AutoMigrate completed successfully")

	return db, nil
}
