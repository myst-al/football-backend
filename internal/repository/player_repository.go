package repository

import (
	"football-backend/internal/models"
	"football-backend/internal/utils"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(p *models.Player) error
	Update(p *models.Player) error
	Delete(id uint) error
	GetAll(q utils.QueryParams) ([]models.Player, int64, error)
	GetByID(id uint) (*models.Player, error)
	GetByTeam(teamID uint) ([]models.Player, error)
	FindJerseyNumber(teamID uint, jerseyNumber int) (*models.Player, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db}
}

func (r *playerRepository) Create(p *models.Player) error {
	return r.db.Create(p).Error
}

func (r *playerRepository) Update(p *models.Player) error {
	return r.db.Save(p).Error
}

func (r *playerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Player{}, id).Error
}

func (r *playerRepository) GetAll(q utils.QueryParams) ([]models.Player, int64, error) {
	var items []models.Player
	var total int64

	db := r.db.Model(&models.Player{}).Preload("Team")

	db = utils.ApplyFilters(db, q)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order(q.Sort + " " + q.Order)

	offset := (q.Page - 1) * q.Limit
	if err := db.Offset(offset).Limit(q.Limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *playerRepository) GetByID(id uint) (*models.Player, error) {
	var p models.Player
	err := r.db.Preload("Team").First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *playerRepository) GetByTeam(teamID uint) ([]models.Player, error) {
	var players []models.Player

	err := r.db.
		Preload("Team").
		Where("team_id = ?", teamID).
		Find(&players).Error

	return players, err
}

func (r *playerRepository) FindJerseyNumber(teamID uint, jerseyNumber int) (*models.Player, error) {
	var p models.Player
	err := r.db.Where("team_id = ? AND jersey_number = ?", teamID, jerseyNumber).First(&p).Error
	return &p, err
}
