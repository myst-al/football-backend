package repository

import (
	"football-backend/internal/models"
	"football-backend/internal/utils"

	"gorm.io/gorm"
)

type TeamRepository interface {
	Create(team *models.Team) error
	GetAll(q utils.QueryParams) ([]models.Team, int64, error)
	GetByID(id uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db}
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *teamRepository) GetAll(q utils.QueryParams) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	db := r.db.Model(&models.Team{})

	db = utils.ApplyFilters(db, q)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order(q.Sort + " " + q.Order)

	offset := (q.Page - 1) * q.Limit
	if err := db.Offset(offset).Limit(q.Limit).Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func (r *teamRepository) GetByID(id uint) (*models.Team, error) {
	var team models.Team

	if err := r.db.First(&team, id).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *teamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *teamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}
