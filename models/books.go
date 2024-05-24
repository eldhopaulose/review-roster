package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Books struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Title     string
    Author    string
    Quantity  int
}

func (b *Books) TableName() string {
    return "books"
}
