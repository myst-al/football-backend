package models

import (
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	MatchID uint
	Match   Match `gorm:"foreignKey:MatchID"`

	TeamID uint
	Team   Team `gorm:"foreignKey:TeamID"`

	ScorerPlayerID uint
	Scorer         Player `gorm:"foreignKey:ScorerPlayerID"`

	Minute string
}
