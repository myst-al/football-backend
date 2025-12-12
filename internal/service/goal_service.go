package service

import (
	"football-backend/internal/dto"
	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"
	"regexp"
)

type GoalService interface {
	AddGoal(g *models.Goal) error
	GetGoals(matchID uint) ([]models.Goal, error)
	TopScorers(limit int) ([]dto.TopScorerDTO, error)
}

type goalService struct {
	repo      repository.GoalRepository
	matchRepo repository.MatchRepository
}

func NewGoalService(
	goalRepo repository.GoalRepository,
	matchRepo repository.MatchRepository,
) GoalService {
	return &goalService{
		repo:      goalRepo,
		matchRepo: matchRepo,
	}
}

func (s *goalService) AddGoal(g *models.Goal) error {
	if g.ScorerPlayerID == 0 {
		return apperror.NewValidationError("pencetak gol wajib diisi")
	}
	if g.MatchID == 0 {
		return apperror.NewValidationError("match_id wajib diisi")
	}
	if g.TeamID == 0 {
		return apperror.NewValidationError("team_id wajib diisi")
	}

	match, err := s.matchRepo.GetByID(g.MatchID)
	if err != nil {
		return apperror.NewNotFoundError("pertandingan tidak ditemukan")
	}

	if g.TeamID != match.HomeTeamID && g.TeamID != match.AwayTeamID {
		return apperror.NewValidationError("team pencetak gol tidak sesuai dengan tim yang bertanding")
	}

	validMinute, _ := regexp.MatchString(`^([0-9]|[1-8][0-9]|90)$|(^(45|90)\+[0-9]+$)`, g.Minute)
	if !validMinute {
		return apperror.NewValidationError(
			"format menit tidak valid. Gunakan angka 0–90 atau hanya 45+X / 90+X (contoh: 45+2, 90+3)",
		)
	}

	invalidStatuses := map[string]bool{
		"DIJADWALKAN": true,
		"SELESAI":     true,
		"DIBATALKAN":  true,
	}

	if invalidStatuses[match.Status] {
		return apperror.NewValidationError(
			"tidak dapat menambahkan gol karena status pertandingan: " + match.Status,
		)
	}

	if err := s.repo.AddGoal(g); err != nil {
		return apperror.NewInternalError("gagal menambahkan gol")
	}

	return nil
}

func (s *goalService) GetGoals(matchID uint) ([]models.Goal, error) {
	goals, err := s.repo.GetGoals(matchID)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data gol")
	}
	return goals, nil
}

func (s *goalService) TopScorers(limit int) ([]dto.TopScorerDTO, error) {
	list, err := s.repo.TopScorers(limit)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil top scorer")
	}
	return list, nil
}
