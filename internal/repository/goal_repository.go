package repository

import (
	"football-backend/internal/dto"
	"football-backend/internal/models"

	"gorm.io/gorm"
)

type GoalRepository interface {
	AddGoal(g *models.Goal) error
	GetGoals(matchID uint) ([]models.Goal, error)
	TopScorers(limit int) ([]dto.TopScorerDTO, error)
}

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{db}
}

func (r *goalRepository) AddGoal(g *models.Goal) error {
	return r.db.Create(g).Error
}

func (r *goalRepository) GetGoals(matchID uint) ([]models.Goal, error) {
	var goals []models.Goal

	err := r.db.
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Preload("Team").
		Preload("Scorer").
		Preload("Scorer.Team").
		Where("match_id = ?", matchID).
		Find(&goals).Error

	return goals, err
}

func (r *goalRepository) TopScorers(limit int) ([]dto.TopScorerDTO, error) {
	var result []dto.TopScorerDTO

	err := r.db.Table("goals AS g").
		Select(`
			g.scorer_player_id AS player_id,
			p.name AS player_name,
			p.team_id AS team_id,
			t.name AS team_name,
			COUNT(*) AS goals
		`).
		Joins("LEFT JOIN players p ON p.id = g.scorer_player_id").
		Joins("LEFT JOIN teams t ON t.id = p.team_id").
		Group("g.scorer_player_id").
		Order("goals DESC").
		Limit(limit).
		Scan(&result).Error

	return result, err
}
