package models

import (
	"api/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		utils.GetEnv("DB_HOST", "127.0.0.1"),
		utils.GetEnv("DB_USER", "swipcord"),
		utils.GetEnv("DB_PASSWORD", "swipcord"),
		utils.GetEnv("DB_NAME", "swipcord"),
		utils.GetEnv("DB_PORT", "5432"))
	var err error = nil

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	return Database, err

}
