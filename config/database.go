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
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_SSLMODE"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db

    log.Println("Running database migrations...")
    
    // ✅ Core System models ONLY
    err = db.AutoMigrate(
        &models.User{},
        &models.AuthLog{},
        &models.Anime{},   
        &models.Episode{}, 
    )
    
    if err != nil {
        log.Fatal("❌ Failed to run migrations:", err)
    }

    log.Println("✓ PostgreSQL connected & migrated successfully")
}
