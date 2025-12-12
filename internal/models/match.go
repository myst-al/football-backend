package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	MatchDateTime time.Time

	HomeTeamID uint
	AwayTeamID uint

	HomeTeam Team `gorm:"foreignKey:HomeTeamID"`
	AwayTeam Team `gorm:"foreignKey:AwayTeamID"`

	Status string `gorm:"type:ENUM('DIJADWALKAN','SEDANG BERLANGSUNG','SELESAI','DIBATALKAN');default:'DIJADWALKAN'"`

	Goals []Goal `gorm:"foreignKey:MatchID"`
}
