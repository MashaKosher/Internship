package db

import (
	"adminservice/internal/config"
	"adminservice/internal/entity"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConncetDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.AppConfig.DB.Host, config.AppConfig.DB.User, config.AppConfig.DB.Password, config.AppConfig.DB.Name, config.AppConfig.DB.Port, config.AppConfig.DB.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed connection to DB: " + err.Error())
	}
	log.Println("Successfully DB connection")
	DB = db
	AutoMigrations(DB)
}

func AutoMigrations(conncetion *gorm.DB) {
	conncetion.AutoMigrate(
		&entity.GameSettings{},
		&entity.Season{},
		&entity.DBDailyTasks{},
	)
}
