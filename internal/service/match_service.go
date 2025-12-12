package service

import (
	"sort"
	"strings"

	"football-backend/internal/dto"
	apperror "football-backend/internal/errors"
	"football-backend/internal/models"
	"football-backend/internal/repository"
	"football-backend/internal/utils"
)

type MatchService interface {
	Create(m *models.Match) error
	Update(m *models.Match) error
	GetByID(id uint) (*models.Match, error)
	GetList(q utils.QueryParams) (map[string]interface{}, error)
	ProcessResult(matchID uint) error
	Report(matchID uint) (map[string]interface{}, error)
	LeagueStanding() ([]dto.StandingDTO, error)
}

type matchService struct {
	repo     repository.MatchRepository
	goalRepo repository.GoalRepository
	teamRepo repository.TeamRepository
}

func NewMatchService(r repository.MatchRepository, g repository.GoalRepository, t repository.TeamRepository) MatchService {
	return &matchService{repo: r, goalRepo: g, teamRepo: t}
}

func (s *matchService) Create(m *models.Match) error {
	if m.HomeTeamID == 0 || m.AwayTeamID == 0 {
		return apperror.NewValidationError("home_team_id dan away_team_id wajib diisi")
	}
	if m.HomeTeamID == m.AwayTeamID {
		return apperror.NewValidationError("home dan away team tidak boleh sama")
	}

	conflict, err := s.repo.CheckConflict(m.HomeTeamID, m.MatchDateTime)
	if err != nil {
		return apperror.NewInternalError("gagal memeriksa jadwal")
	}
	if conflict {
		return apperror.NewConflictError("jadwal bentrok dengan pertandingan lain (tim home)")
	}

	conflict, err = s.repo.CheckConflict(m.AwayTeamID, m.MatchDateTime)
	if err != nil {
		return apperror.NewInternalError("gagal memeriksa jadwal")
	}
	if conflict {
		return apperror.NewConflictError("jadwal bentrok dengan pertandingan lain (tim away)")
	}

	if err := s.repo.Create(m); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperror.NewConflictError("jadwal pertandingan sudah ada")
		}
		return apperror.NewInternalError("gagal membuat pertandingan")
	}
	return nil
}

func (s *matchService) Update(m *models.Match) error {
	if err := s.repo.Update(m); err != nil {
		return apperror.NewInternalError("gagal memperbarui pertandingan")
	}
	return nil
}

func (s *matchService) GetByID(id uint) (*models.Match, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("pertandingan tidak ditemukan")
	}
	return m, nil
}

func (s *matchService) GetList(q utils.QueryParams) (map[string]interface{}, error) {
	items, total, err := s.repo.GetAll(q)
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil data pertandingan")
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

func (s *matchService) ProcessResult(matchID uint) error {
	match, err := s.repo.GetByID(matchID)
	if err != nil {
		return apperror.NewNotFoundError("pertandingan tidak ditemukan")
	}
	goals, err := s.goalRepo.GetGoals(matchID)
	if err != nil {
		return apperror.NewInternalError("gagal mengambil data gol")
	}

	homeScore := 0
	awayScore := 0

	for _, g := range goals {
		if g.TeamID == match.HomeTeamID {
			homeScore++
		} else {
			awayScore++
		}
	}

	if homeScore > awayScore {
		match.Status = "HOME_WIN"
	} else if awayScore > homeScore {
		match.Status = "AWAY_WIN"
	} else {
		match.Status = "DRAW"
	}

	if err := s.repo.Update(match); err != nil {
		return apperror.NewInternalError("gagal menyimpan hasil pertandingan")
	}
	return nil
}

func (s *matchService) Report(matchID uint) (map[string]interface{}, error) {
	match, err := s.repo.GetByID(matchID)
	if err != nil {
		return nil, apperror.NewNotFoundError("match tidak ditemukan")
	}

	goals, _ := s.goalRepo.GetGoals(matchID)

	homeScore := 0
	awayScore := 0
	scorerCount := map[uint]int{}
	scorerDetail := map[uint]models.Player{}

	goalList := []map[string]interface{}{}

	for _, g := range goals {
		if g.TeamID == match.HomeTeamID {
			homeScore++
		} else {
			awayScore++
		}
		scorerCount[g.ScorerPlayerID]++
		scorerDetail[g.ScorerPlayerID] = g.Scorer

		goalList = append(goalList, map[string]interface{}{
			"player": map[string]interface{}{
				"id":            g.Scorer.ID,
				"name":          g.Scorer.Name,
				"position":      g.Scorer.Position,
				"jersey_number": g.Scorer.JerseyNumber,
			},
			"minute": g.Minute,
		})
	}

	topScorer := map[string]interface{}{}
	topGoals := 0

	for playerID, total := range scorerCount {
		if total > topGoals {
			topGoals = total
			p := scorerDetail[playerID]

			topScorer = map[string]interface{}{
				"id":            p.ID,
				"name":          p.Name,
				"height":        p.HeightCM,
				"weight":        p.WeightKG,
				"position":      p.Position,
				"jersey_number": p.JerseyNumber,
				"goals":         total,
			}
		}
	}

	homeWins, _ := s.repo.CountHomeWins(match.HomeTeamID)
	awayWins, _ := s.repo.CountAwayWins(match.AwayTeamID)

	finalStatus := ""
	if homeScore > awayScore {
		finalStatus = "Tim Home Menang"
	} else if awayScore > homeScore {
		finalStatus = "Tim Away Menang"
	} else {
		finalStatus = "Draw"
	}

	return map[string]interface{}{
		"match": map[string]interface{}{
			"id":         match.ID,
			"match_date": match.MatchDateTime,
			"home_team":  match.HomeTeam,
			"away_team":  match.AwayTeam,
		},
		"score": map[string]interface{}{
			"home": homeScore,
			"away": awayScore,
		},
		"status":     finalStatus,
		"goals":      goalList,
		"top_scorer": topScorer,
		"home_wins":  homeWins,
		"away_wins":  awayWins,
	}, nil
}

func (s *matchService) LeagueStanding() ([]dto.StandingDTO, error) {
	teams, _, err := s.teamRepo.GetAll(utils.QueryParams{
		Page:    1,
		Limit:   9999,
		Sort:    "id",
		Order:   "ASC",
		Filters: map[string]map[utils.FilterOperator]string{},
	})
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil daftar tim")
	}

	matches, err := s.repo.GetFinishedMatches()
	if err != nil {
		return nil, apperror.NewInternalError("gagal mengambil pertandingan selesai")
	}

	standing := map[uint]*dto.StandingDTO{}

	for _, t := range teams {
		standing[t.ID] = &dto.StandingDTO{
			TeamID:   t.ID,
			TeamName: t.Name,
		}
	}

	for _, m := range matches {
		home := standing[m.HomeTeamID]
		away := standing[m.AwayTeamID]

		homeGoals := 0
		awayGoals := 0

		for _, g := range m.Goals {
			if g.TeamID == m.HomeTeamID {
				homeGoals++
			} else {
				awayGoals++
			}
		}

		home.Played++
		away.Played++

		home.GoalsFor += homeGoals
		home.GoalsAgainst += awayGoals

		away.GoalsFor += awayGoals
		away.GoalsAgainst += homeGoals

		if homeGoals > awayGoals {
			home.Wins++
			away.Losses++
		} else if awayGoals > homeGoals {
			away.Wins++
			home.Losses++
		} else {
			home.Draws++
			away.Draws++
		}
	}

	list := []dto.StandingDTO{}
	for _, s := range standing {
		s.GoalDifference = s.GoalsFor - s.GoalsAgainst
		s.Points = s.Wins*3 + s.Draws
		list = append(list, *s)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Points != list[j].Points {
			return list[i].Points > list[j].Points
		}
		if list[i].GoalDifference != list[j].GoalDifference {
			return list[i].GoalDifference > list[j].GoalDifference
		}
		return list[i].GoalsFor > list[j].GoalsFor
	})

	return list, nil
}
