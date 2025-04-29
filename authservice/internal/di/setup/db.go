package setup

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func mustDB(cfg di.ConfigType, logger di.LoggerType) di.DBType {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatal("Failed connection to DB: " + err.Error())
	}
	logger.Info("Successful connection to DB")
	go autoMigrations(db)
	return db
}

func autoMigrations(conncetion *gorm.DB) {
	conncetion.AutoMigrate(
		&entity.User{},
	)
}
