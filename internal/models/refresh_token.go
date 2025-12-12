package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID uint   `gorm:"index"`
	Token  string `gorm:"size:512;unique;not null"`
	JTI    string `gorm:"size:128"`

	ExpiresAt time.Time
}
