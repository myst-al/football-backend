package repository

import (
	"football-backend/internal/models"

	"gorm.io/gorm"
)

type PlayerTransferRepository interface {
	Create(t *models.PlayerTransfer) error
}

type playerTransferRepository struct {
	db *gorm.DB
}

func NewPlayerTransferRepository(db *gorm.DB) PlayerTransferRepository {
	return &playerTransferRepository{db}
}

func (r *playerTransferRepository) Create(t *models.PlayerTransfer) error {
	return r.db.Create(t).Error
}
