package config

import (
	"fmt"
	"log"
	"os"

	"core-anime/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Debug: Print DSN (remove password for security)
	log.Printf("Connecting with DSN: host=%s user=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Assign to global DB
	DB = db

	// Verify connection
	var currentDB string
	err = db.Raw("SELECT current_database()").Scan(&currentDB).Error
	if err != nil {
		log.Fatal("Failed to verify database connection:", err)
	}
	log.Printf("✓ Connected to database: %s", currentDB)

	// Check current table count BEFORE migration
	var tableCountBefore int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCountBefore)
	log.Printf("Tables in public schema (before migration): %d", tableCountBefore)

	// Run migrations with error checking
	log.Println("Running database migrations...")
	
	err = db.AutoMigrate(
		&models.User{},
		&models.Genre{},
		&models.Anime{},
		&models.Episode{},
	)
	
	if err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	// Check table count AFTER migration
	var tableCountAfter int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCountAfter)
	log.Printf("Tables in public schema (after migration): %d", tableCountAfter)

	// List all tables created
	var tables []string
	db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)
	
	if len(tables) > 0 {
		log.Println("✓ Tables created:")
		for _, table := range tables {
			log.Printf("  - %s", table)
		}
	} else {
		log.Println("⚠️  WARNING: No tables found after migration!")
	}

	log.Println("✓ PostgreSQL connected & migrated successfully")
}
