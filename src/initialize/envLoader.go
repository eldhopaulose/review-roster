package initialize

import (
    "os"
    "github.com/joho/godotenv"
)

func LoadEnv() {
    if os.Getenv("ENV") != "PRODUCTION" {
        err := godotenv.Load(".env")
        if err != nil {
            panic("Error loading .env file")
        }
    }
}
