package service

import (
	"strings"

	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"
	"football-backend/internal/utils"
)

type PlayerService interface {
	Create(p *models.Player) error
	Update(p *models.Player) error
	Delete(id uint) error
	GetList(q utils.QueryParams) (map[string]interface{}, error)
	GetByID(id uint) (*models.Player, error)
	GetByTeam(teamID uint) ([]models.Player, error)
	TransferPlayer(playerID, newTeamID uint, newJersey int) error
}

type playerService struct {
	repo         repository.PlayerRepository
	transferRepo repository.PlayerTransferRepository
	teamRepo     repository.TeamRepository
}

func NewPlayerService(repo repository.PlayerRepository, tRepo repository.PlayerTransferRepository, teamRepo repository.TeamRepository) PlayerService {
	return &playerService{repo: repo, transferRepo: tRepo, teamRepo: teamRepo}
}

func validatePosition(pos string) bool {
	valid := map[string]bool{
		"PENYERANG":      true,
		"GELANDANG":      true,
		"BERTAHAN":       true,
		"PENJAGA_GAWANG": true,
	}
	return valid[pos]
}

func (s *playerService) GetList(q utils.QueryParams) (map[string]interface{}, error) {
	items, total, err := s.repo.GetAll(q)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data pemain")
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

func (s *playerService) Create(p *models.Player) error {
	if p.Name == "" {
		return apperror.NewValidationError("nama pemain wajib diisi")
	}

	if !validatePosition(p.Position) {
		return apperror.NewValidationError("posisi pemain tidak valid")
	}

	team, err := s.teamRepo.GetByID(p.TeamID)
	if err != nil {
		return apperror.NewNotFoundError("team tidak ditemukan atau sudah dihapus")
	}
	if team.DeletedAt.Valid {
		return apperror.NewValidationError("tidak bisa menambahkan pemain ke team yang sudah dihapus")
	}

	if _, err := s.repo.FindJerseyNumber(p.TeamID, p.JerseyNumber); err == nil {
		return apperror.NewConflictError("nomor punggung sudah digunakan dalam tim ini")
	}

	if err := s.repo.Create(p); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperror.NewConflictError("nomor punggung sudah digunakan")
		}
		return apperror.NewInternalError("gagal membuat pemain")
	}

	saved, err := s.repo.GetByID(p.ID)
	if err != nil {
		return apperror.NewInternalError("gagal memuat data pemain")
	}
	*p = *saved

	return nil
}

func (s *playerService) Update(p *models.Player) error {
	if !validatePosition(p.Position) {
		return apperror.NewValidationError("posisi pemain tidak valid")
	}

	exist, err := s.repo.FindJerseyNumber(p.TeamID, p.JerseyNumber)
	if err == nil && exist != nil && exist.ID != p.ID {
		return apperror.NewConflictError("nomor punggung sudah digunakan pemain lain")
	}

	if err := s.repo.Update(p); err != nil {
		return apperror.NewInternalError("gagal memperbarui pemain")
	}
	return nil
}

func (s *playerService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return apperror.NewNotFoundError("pemain tidak ditemukan")
	}
	if err := s.repo.Delete(id); err != nil {
		return apperror.NewInternalError("gagal menghapus pemain")
	}
	return nil
}

func (s *playerService) GetByID(id uint) (*models.Player, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("pemain tidak ditemukan")
	}
	return p, nil
}

func (s *playerService) GetByTeam(teamID uint) ([]models.Player, error) {
	list, err := s.repo.GetByTeam(teamID)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data pemain")
	}
	return list, nil
}

func (s *playerService) TransferPlayer(playerID, newTeamID uint, newJersey int) error {
	player, err := s.repo.GetByID(playerID)
	if err != nil {
		return apperror.NewNotFoundError("pemain tidak ditemukan")
	}

	newTeam, err := s.teamRepo.GetByID(newTeamID)
	if err != nil {
		return apperror.NewNotFoundError("team tujuan tidak ditemukan")
	}
	if newTeam.DeletedAt.Valid {
		return apperror.NewValidationError("tidak bisa mentransfer pemain ke team yang sudah dihapus")
	}

	exist, _ := s.repo.FindJerseyNumber(newTeamID, newJersey)
	if exist != nil {
		return apperror.NewConflictError("nomor punggung sudah dipakai di tim baru")
	}

	transfer := &models.PlayerTransfer{
		PlayerID:     player.ID,
		OldTeamID:    player.TeamID,
		NewTeamID:    newTeamID,
		JerseyNumber: newJersey,
	}
	if err := s.transferRepo.Create(transfer); err != nil {
		return apperror.NewInternalError("gagal menyimpan riwayat transfer")
	}

	player.TeamID = newTeamID
	player.JerseyNumber = newJersey

	if err := s.repo.Update(player); err != nil {
		return apperror.NewInternalError("gagal memperbarui data pemain")
	}
	return nil
}
