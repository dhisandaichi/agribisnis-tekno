package db

import (
	_ "embed"
	"fmt"
	"log"
	"sync"
	"time"

	"agritec/backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed schema.sql
var Schema string

var (
	once    sync.Once
	conn    *gorm.DB
	connErr error
)

// DSN: pakai DATABASE_URL (Neon/Supabase/Vercel Postgres) bila ada, jika tidak susun dari DB_*.
func dsn() string {
	if url := config.GetEnv("DATABASE_URL", config.GetEnv("POSTGRES_URL", "")); url != "" {
		return url
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", "agritec_user"),
		config.GetEnv("DB_PASSWORD", "agritec_password"),
		config.GetEnv("DB_NAME", "agritec_db"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_SSLMODE", "require"),
	)
}

// Get mengembalikan koneksi tunggal (dipakai ulang antar request serverless).
func Get() (*gorm.DB, error) {
	once.Do(func() {
		conn, connErr = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn(),
			PreferSimpleProtocol: true, // aman untuk connection pooler (PgBouncer / Neon pooler)
		}), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
		if connErr != nil {
			log.Println("database connection failed:", connErr)
			return
		}
		sqlDB, err := conn.DB()
		if err != nil {
			connErr = err
			return
		}
		sqlDB.SetMaxOpenConns(5)
		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	})
	return conn, connErr
}

// Migrate menjalankan schema.sql (idempotent).
func Migrate(g *gorm.DB) error {
	return g.Exec(Schema).Error
}
