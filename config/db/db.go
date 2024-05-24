package db

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/eldhopaulose/ReviewRoster/models"
)

var DB *gorm.DB

func ConnectDB() {
    var err error
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSL_MODE"),
        os.Getenv("DB_PASSWORD"),
    )
    fmt.Println("Connecting to database with DSN:", dsn) // Add this line for debugging
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic(fmt.Sprintf("Failed to connect to the database, got error %v", err))
    }
}

func GetDB() *gorm.DB {
    return DB
}

func Migrate() {
    err := DB.AutoMigrate(&models.Books{})
    if err != nil {
        panic("Failed to migrate database!")
    }
}
