package service

import (
	"strings"

	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"
	"football-backend/internal/utils"
)

type TeamService interface {
	Create(team *models.Team) error
	GetList(q utils.QueryParams) (map[string]interface{}, error)
	GetByID(id uint) (*models.Team, error)
	Update(team *models.Team) error
	Delete(id uint) error
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(r repository.TeamRepository) TeamService {
	return &teamService{repo: r}
}

func (s *teamService) Create(team *models.Team) error {
	if team.Name == "" {
		return apperror.NewValidationError("nama team wajib diisi")
	}

	if err := s.repo.Create(team); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperror.NewConflictError("nama team sudah digunakan")
		}
		return apperror.NewInternalError("gagal membuat team")
	}
	return nil
}

func (s *teamService) GetList(q utils.QueryParams) (map[string]interface{}, error) {
	items, total, err := s.repo.GetAll(q)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data tim")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}

	return map[string]interface{}{
		"items": items,
		"pagination": map[string]interface{}{
			"page":        q.Page,
			"limit":       q.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}, nil
}

func (s *teamService) GetByID(id uint) (*models.Team, error) {
	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("team tidak ditemukan")
	}
	return team, nil
}

func (s *teamService) Update(team *models.Team) error {
	if team.Name == "" {
		return apperror.NewValidationError("nama team wajib diisi")
	}
	if err := s.repo.Update(team); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperror.NewConflictError("nama team sudah digunakan")
		}
		return apperror.NewInternalError("gagal memperbarui team")
	}
	return nil
}

func (s *teamService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.NewNotFoundError("team tidak ditemukan")
	}
	if err := s.repo.Delete(id); err != nil {
		return apperror.NewInternalError("gagal menghapus team")
	}
	return nil
}
