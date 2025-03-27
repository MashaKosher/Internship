package db

import (
	"authservice/internal/config"
	"authservice/internal/db/models"
	"authservice/internal/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConncetDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Cfg.DB.Host, config.Cfg.DB.User, config.Cfg.DB.Password, config.Cfg.DB.Name, config.Cfg.DB.Port, config.Cfg.DB.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Error("Failed connection to DB: " + err.Error())
		panic("Failed connection to DB: " + err.Error())
	}
	logger.Logger.Info("Successful connection to DB")
	DB = db
	AutoMigrations(DB)
}

func AutoMigrations(conncetion *gorm.DB) {
	conncetion.Debug().AutoMigrate(
		&models.User{},
	)
}
