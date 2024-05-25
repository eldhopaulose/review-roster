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
    Title     string         `gorm:"not null"`
    Author    string         `gorm:"not null"`
    Genre     string         `gorm:"not null;default:'Unknown'"` // Default value added
    Year      int            `gorm:"not null;default:0"`          // Default value added
    Rating    int            `gorm:"not null;default:0"`          // Default value added
    Review    string         `gorm:"not null;default:''"`         // Default value added
}

func (b *Books) TableName() string {
    return "books"
}
