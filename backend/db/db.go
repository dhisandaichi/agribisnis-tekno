package db

import (
	_ "embed"
	"fmt"
	"log"
	"net/url"
	"strings"
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
	mu   sync.Mutex
	conn *gorm.DB
)

// cleanURL merapikan connection string yang ditempel dari dashboard (spasi, kutip,
// awalan "psql") dan membuang parameter yang tidak dikenal driver pgx, misalnya
// channel_binding dari Neon, yang jika ikut terkirim membuat koneksi ditolak server.
func cleanURL(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "psql ")
	s = strings.Trim(s, "'\" ")
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") {
		return s
	}
	q := u.Query()
	q.Del("channel_binding")
	if q.Get("sslmode") == "" {
		q.Set("sslmode", "require")
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// DSN: pakai DATABASE_URL (Neon/Supabase/Vercel Postgres) bila ada, jika tidak susun dari DB_*.
func dsn() string {
	if u := config.GetEnv("DATABASE_URL", config.GetEnv("POSTGRES_URL", "")); u != "" {
		return cleanURL(u)
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
// Kegagalan tidak disimpan: request berikutnya mencoba lagi (mis. Neon sedang bangun).
func Get() (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		return conn, nil
	}
	var g *gorm.DB
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		g, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn(),
			PreferSimpleProtocol: true, // aman untuk connection pooler (PgBouncer / Neon pooler)
		}), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
		if err == nil {
			break
		}
		log.Printf("database connection failed (percobaan %d): %v", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := g.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	conn = g
	return conn, nil
}

// Migrate menjalankan schema.sql (idempotent).
func Migrate(g *gorm.DB) error {
	return g.Exec(Schema).Error
}
