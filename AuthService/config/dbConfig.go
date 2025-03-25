package config

import (
	"fmt"

	models "authservice/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConncetDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", Envs.DBHost, Envs.DBUser, Envs.DBPassword, Envs.DBName, Envs.DBPort, Envs.DBSSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Error("Failed connection to DB")
	}
	Logger.Info("Successful connection to DB")
	DB = db
	AutoMigrations(DB)
}

func AutoMigrations(conncetion *gorm.DB) {
	conncetion.Debug().AutoMigrate(
		&models.User{},
	)
}
