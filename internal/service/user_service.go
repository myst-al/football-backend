package service

import (
	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"
)

type UserService interface {
	GetAdmins() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAdmins() ([]models.User, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data admin")
	}
	return list, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("admin tidak ditemukan")
	}
	return u, nil
}

func (s *userService) Delete(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.NewNotFoundError("admin tidak ditemukan")
	}

	if user.ID != id {
		return apperror.NewValidationError("tidak dapat menghapus diri sendiri")
	}

	if err := s.repo.Delete(id); err != nil {
		return apperror.NewInternalError("gagal menghapus admin")
	}
	return nil
}
