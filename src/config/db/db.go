package db

import (
    "fmt"
    "os"
    "database/sql"
    "github.com/lib/pq"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/eldhopaulose/ReviewRoster/src/models"
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

    // Open a new connection to the PostgreSQL server without specifying a database.
    conn, err := sql.Open("postgres", fmt.Sprintf(
        "host=%s port=%s user=%s sslmode=%s password=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_SSL_MODE"),
        os.Getenv("DB_PASSWORD"),
    ))
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to the database server, got error %v", err))
    }
    defer conn.Close()

    // Create the database if it doesn't exist
    _, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", os.Getenv("DB_NAME")))
    if err != nil {
        // If the error is related to the database already existing, ignore it
        if e, ok := err.(*pq.Error); ok && e.Code == "42P04" {
            fmt.Println("Database already exists")
        } else {
            panic(fmt.Sprintf("Failed to create database, got error %v", err))
        }
    } else {
        fmt.Println("Database created successfully")
    }

    // Connect to the created database
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
