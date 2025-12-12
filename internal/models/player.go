package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TeamID uint `json:"team_id"`
	Team   Team `gorm:"foreignKey:TeamID"`

	Name         string `gorm:"size:255;not null"`
	HeightCM     int    `json:"height"`
	WeightKG     int    `json:"weight"`
	Position     string `gorm:"type:ENUM('PENYERANG','GELANDANG','BERTAHAN','PENJAGA_GAWANG');not null"`
	JerseyNumber int    `gorm:"not null"`
}
