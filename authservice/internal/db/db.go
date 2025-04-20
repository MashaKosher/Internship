package db

import (
	"authservice/internal/config"
	"authservice/internal/entity"
	"authservice/internal/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConncetDB() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.AppConfig.DB.Host, config.AppConfig.DB.User, config.AppConfig.DB.Password, config.AppConfig.DB.Name, config.AppConfig.DB.Port, config.AppConfig.DB.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Logger.Fatal("Failed connection to DB: " + err.Error())
	}
	logger.Logger.Info("Successful connection to DB")
	DB = db
	AutoMigrations(DB)

}

func AutoMigrations(conncetion *gorm.DB) {
	conncetion.AutoMigrate(
		&entity.User{},
	)
}
