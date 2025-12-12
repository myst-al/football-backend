package dto

import (
	"football-backend/internal/models"
	"time"
)

type MatchDTO struct {
	ID        uint      `json:"id"`
	MatchDate string    `json:"match_date"`
	Status    string    `json:"status"`
	HomeTeam  TeamDTO   `json:"home_team"`
	AwayTeam  TeamDTO   `json:"away_team"`
	Goals     []GoalDTO `json:"goals"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
}

func ToMatchDTO(m *models.Match) MatchDTO {
	homeScore := 0
	awayScore := 0
	for _, g := range m.Goals {
		if g.TeamID == m.HomeTeamID {
			homeScore++
		} else {
			awayScore++
		}
	}

	goals := make([]GoalDTO, 0)
	for _, g := range m.Goals {
		goals = append(goals, ToGoalDTO(&g))
	}

	return MatchDTO{
		ID:        m.ID,
		MatchDate: m.MatchDateTime.Format(time.RFC3339),
		Status:    m.Status,
		HomeTeam:  ToTeamDTO(&m.HomeTeam),
		AwayTeam:  ToTeamDTO(&m.AwayTeam),
		Goals:     goals,
		HomeScore: homeScore,
		AwayScore: awayScore,
	}
}

func ToMatchDTOList(list []models.Match) []MatchDTO {
	result := make([]MatchDTO, 0)
	for _, m := range list {
		result = append(result, ToMatchDTO(&m))
	}
	return result
}
