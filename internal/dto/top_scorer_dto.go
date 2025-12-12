package dto

type TopScorerDTO struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamID     uint   `json:"team_id"`
	TeamName   string `json:"team_name"`
	Goals      int    `json:"goals"`
}
