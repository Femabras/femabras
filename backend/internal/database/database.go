// femabras/backend/internal/database/database.go
package database

import (
	"log"
	"time"

	"github.com/Femabras/femabras/backend/internal/models"

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

	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Running AutoMigrate...")

	err = db.AutoMigrate(
		&models.PendingUser{},
		&models.Challenge{},
		&models.UserStat{},
		&models.User{},
		&models.OTP{},
	)
	if err != nil {
		return nil, err
	}
	log.Println("AutoMigrate completed successfully")

	return db, nil
}
