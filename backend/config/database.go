package config

import (
	"fmt"
	"log"

	"github.com/agritec/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := GetEnv("SUPABASE_DB_URL", "")

	if dsn == "" {
		host := GetEnv("DB_HOST", "localhost")
		user := GetEnv("DB_USER", "agritec_user")
		password := GetEnv("DB_PASSWORD", "agritec_password")
		dbname := GetEnv("DB_NAME", "agritec_db")
		port := GetEnv("DB_PORT", "5432")
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	} else {
		log.Println("Using SUPABASE_DB_URL from environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate all models – creates/updates tables automatically
	err = db.AutoMigrate(
		&models.User{},
		&models.Farmer{},
		&models.Visit{},
		&models.Distributor{},
		&models.Stock{},
		&models.Debt{},
		&models.Content{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	log.Println("Database connected and all tables migrated successfully")
	return db
}
