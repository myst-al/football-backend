package repository

import (
	"football-backend/internal/models"
	"football-backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(m *models.Match) error
	Update(m *models.Match) error
	GetByID(id uint) (*models.Match, error)
	GetAll(q utils.QueryParams) ([]models.Match, int64, error)
	CountHomeWins(teamID uint) (int64, error)
	CountAwayWins(teamID uint) (int64, error)
	CheckConflict(teamID uint, datetime time.Time) (bool, error)
	GetFinishedMatches() ([]models.Match, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db}
}

func (r *matchRepository) Create(m *models.Match) error {
	return r.db.Create(m).Error
}

func (r *matchRepository) Update(m *models.Match) error {
	return r.db.Save(m).Error
}

func (r *matchRepository) GetByID(id uint) (*models.Match, error) {
	var match models.Match

	err := r.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Goals").
		Preload("Goals.Team").
		Preload("Goals.Scorer").
		Preload("Goals.Scorer.Team").
		First(&match, id).Error

	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) GetAll(q utils.QueryParams) ([]models.Match, int64, error) {
	var items []models.Match
	var total int64

	db := r.db.Model(&models.Match{}).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Goals").
		Preload("Goals.Team").
		Preload("Goals.Scorer").
		Preload("Goals.Scorer.Team")

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

func (r *matchRepository) CountHomeWins(teamID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Match{}).
		Where("home_team_id = ? AND status = 'HOME_WIN'", teamID).
		Count(&count).Error
	return count, err
}

func (r *matchRepository) CountAwayWins(teamID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Match{}).
		Where("away_team_id = ? AND status = 'AWAY_WIN'", teamID).
		Count(&count).Error
	return count, err
}

func (r *matchRepository) CheckConflict(teamID uint, datetime time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Match{}).
		Where("(home_team_id = ? OR away_team_id = ?) AND match_date_time = ?", teamID, teamID, datetime).
		Count(&count).Error

	return count > 0, err
}

func (r *matchRepository) GetFinishedMatches() ([]models.Match, error) {
	var matches []models.Match

	err := r.db.
		Preload("Goals").
		Preload("Goals.Team").
		Where("status = ?", "SELESAI").
		Find(&matches).Error

	return matches, err
}
