package models

import (
	"api/utils"
	"context"
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		utils.GetEnv("DB_HOST", "127.0.0.1"),
		utils.GetEnv("DB_USER", "swipcord"),
		utils.GetEnv("DB_PASSWORD", "swipcord"),
		utils.GetEnv("DB_NAME", "swipcord"),
		utils.GetEnv("DB_PORT", "5432"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

type contextKey string

const databaseKey contextKey = "database"

func ContextWithDatabase(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, databaseKey, db)
}

func GetDatabase(r *http.Request) *gorm.DB {
	return r.Context().Value(databaseKey).(*gorm.DB)
}
