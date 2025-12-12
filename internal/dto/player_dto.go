package dto

import "football-backend/internal/models"

type TeamSimpleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PlayerDTO struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	HeightCM     int           `json:"height"`
	WeightKG     int           `json:"weight"`
	Position     string        `json:"position"`
	JerseyNumber int           `json:"jersey_number"`
	Team         TeamSimpleDTO `json:"team"`
}

func ToPlayerDTO(p *models.Player) PlayerDTO {
	return PlayerDTO{
		ID:           p.ID,
		Name:         p.Name,
		HeightCM:     p.HeightCM,
		WeightKG:     p.WeightKG,
		Position:     p.Position,
		JerseyNumber: p.JerseyNumber,
		Team: TeamSimpleDTO{
			ID:   p.Team.ID,
			Name: p.Team.Name,
		},
	}
}

func ToPlayerDTOList(list []models.Player) []PlayerDTO {
	result := make([]PlayerDTO, 0, len(list))
	for _, p := range list {
		result = append(result, ToPlayerDTO(&p))
	}
	return result
}
