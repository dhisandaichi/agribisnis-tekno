package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	// Untuk Supabase / Production, sangat disarankan menggunakan connection string (DATABASE_URL)
	dsn := os.Getenv("DATABASE_URL")
	
	if dsn == "" {
		host := GetEnv("DB_HOST", "localhost")
		user := GetEnv("DB_USER", "agritec_user")
		password := GetEnv("DB_PASSWORD", "agritec_password")
		dbname := GetEnv("DB_NAME", "agritec_db")
		port := GetEnv("DB_PORT", "5432")
		// Pastikan SSL Mode = require untuk Supabase jika menggunakan DSN manual
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	} else {
		// Log agar kita tahu bahwa URL Supabase berhasil tertangkap
		log.Println("Using DATABASE_URL from environment variables (Supabase Ready)")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
