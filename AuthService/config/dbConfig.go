package config

import (
	"fmt"
	"os"

	models "authservice/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConncetDB() {
	godotenv.Load()
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, "disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Fatal("failed to connect to the database:")
	}

	Logger.Info("All is good")

	DB = db

	AutoMigrations(DB)

	// tables, err := DB.Migrator().GetTables()
	// if err != nil {
	// 	Logger.Fatal("failed to get tables:" + err.Error())
	// }

	// // Вывод всех таблиц
	// for _, table := range tables {
	// 	Logger.Info(table)
	// }

}

func AutoMigrations(conncetion *gorm.DB) {
	conncetion.Debug().AutoMigrate(
		&models.User{},
	)
}
