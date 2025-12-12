package dto

import "football-backend/internal/models"

type TeamDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	YearFounded int    `json:"year_founded"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

func ToTeamDTO(t *models.Team) TeamDTO {
	return TeamDTO{
		ID:          t.ID,
		Name:        t.Name,
		LogoURL:     t.LogoURL,
		YearFounded: t.YearFounded,
		Address:     t.Address,
		City:        t.City,
	}
}

func ToTeamDTOList(list []models.Team) []TeamDTO {
	result := make([]TeamDTO, 0, len(list))
	for _, t := range list {
		result = append(result, ToTeamDTO(&t))
	}
	return result
}
