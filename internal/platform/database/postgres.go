package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBStore *gorm.DB

func InitializePostgres(dsn string) (*gorm.DB, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDB(*cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// ===== Connection Pool Settings =====
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// ===== Ping Database =====
	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	DBStore = db

	return db, nil
}

// ===== Database Migration =====
func Migrate(db *gorm.DB) error {
	return nil
}
