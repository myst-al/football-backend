package models

import (
	"time"

	"gorm.io/gorm"
)

type PlayerTransfer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PlayerID     uint `json:"player_id"`
	OldTeamID    uint `json:"old_team_id"`
	NewTeamID    uint `json:"new_team_id"`
	JerseyNumber int  `json:"jersey_number"`
}
