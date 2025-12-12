package dto

import "football-backend/internal/models"

type GoalDTO struct {
	ID     uint      `json:"id"`
	Minute string    `json:"minute"`
	Team   TeamDTO   `json:"team"`
	Scorer PlayerDTO `json:"scorer"`
}

func ToGoalDTO(g *models.Goal) GoalDTO {
	return GoalDTO{
		ID:     g.ID,
		Minute: g.Minute,
		Team:   ToTeamDTO(&g.Team),
		Scorer: ToPlayerDTO(&g.Scorer),
	}
}
