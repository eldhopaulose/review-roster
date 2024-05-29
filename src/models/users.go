package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Users struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"not null"`
	Email     string         `gorm:"not null;unique"`
	Password  string         `gorm:"not null"`

}


func (u *Users) TableName() string {
	return "users"
}

