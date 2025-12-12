package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"size:255;unique;not null"`
	LogoURL     string `gorm:"size:1024"`
	YearFounded int
	Address     string `gorm:"size:1024"`
	City        string `gorm:"size:255"`

	Players []Player `gorm:"foreignKey:TeamID"`

	HomeMatches []Match `gorm:"foreignKey:HomeTeamID"`
	AwayMatches []Match `gorm:"foreignKey:AwayTeamID"`
}
