package db

import (
	"authservice/internal/config"
	"authservice/internal/entity"
	"authservice/internal/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConncetDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.AppConfig.DB.Host, config.AppConfig.DB.User, config.AppConfig.DB.Password, config.AppConfig.DB.Name, config.AppConfig.DB.Port, config.AppConfig.DB.SSLMode)
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
		&entity.User{},
	)
}
