package handler

import (
	"football-backend/internal/dto"
	"football-backend/internal/models"
	"football-backend/internal/response"
	"football-backend/internal/service"
	"football-backend/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service     service.MatchService
	goalService service.GoalService
}

func NewMatchHandler(s service.MatchService, g service.GoalService) *MatchHandler {
	return &MatchHandler{s, g}
}

func (h *MatchHandler) Create(c *gin.Context) {
	var input struct {
		MatchDateTime string `json:"match_date_time" binding:"required"`
		HomeTeamID    uint   `json:"home_team_id" binding:"required"`
		AwayTeamID    uint   `json:"away_team_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	dt, err := time.Parse(time.RFC3339, input.MatchDateTime)
	if err != nil {
		response.Error(c, 400, "Format tanggal harus RFC3339")
		return
	}

	m := models.Match{
		MatchDateTime: dt,
		HomeTeamID:    input.HomeTeamID,
		AwayTeamID:    input.AwayTeamID,
	}

	if err := h.service.Create(&m); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 201, "Pertandingan berhasil dibuat", nil)
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	m, err := h.service.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Data pertandingan berhasil diambil", dto.ToMatchDTO(m))
}

func (h *MatchHandler) GetAll(c *gin.Context) {
	q := utils.ParseQuery(c)

	result, err := h.service.GetList(q)
	if err != nil {
		response.FromError(c, err)
		return
	}

	items := result["items"].([]models.Match)

	dtoList := dto.ToMatchDTOList(items)
	result["items"] = dtoList

	response.Success(c, 200, "Data pertandingan berhasil diambil", result)
}

func (h *MatchHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	m, err := h.service.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	var input struct {
		Status        string `json:"status"`
		MatchDateTime string `json:"match_date_time"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	if input.Status != "" {
		m.Status = input.Status
	}

	if input.MatchDateTime != "" {
		if dt, err := time.Parse(time.RFC3339, input.MatchDateTime); err == nil {
			m.MatchDateTime = dt
		}
	}

	if err := h.service.Update(m); err != nil {
		response.FromError(c, err)
		return
	}

	updated, _ := h.service.GetByID(uint(id))
	response.Success(c, 200, "Pertandingan berhasil diperbarui", dto.ToMatchDTO(updated))
}

func (h *MatchHandler) SubmitResult(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Goals []struct {
			TeamID         uint   `json:"team_id"`
			ScorerPlayerID uint   `json:"scorer_player_id"`
			Minute         string `json:"minute"`
		}
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	for _, g := range input.Goals {
		newGoal := models.Goal{
			MatchID:        uint(id),
			TeamID:         g.TeamID,
			ScorerPlayerID: g.ScorerPlayerID,
			Minute:         g.Minute,
		}
		if err := h.goalService.AddGoal(&newGoal); err != nil {
			response.FromError(c, err)
			return
		}
	}

	if err := h.service.ProcessResult(uint(id)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Hasil pertandingan berhasil disimpan", nil)
}

func (h *MatchHandler) Report(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	data, err := h.service.Report(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Report pertandingan berhasil diambil", data)
}

func (h *MatchHandler) Standing(c *gin.Context) {
	data, err := h.service.LeagueStanding()
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Standing berhasil diambil", data)
}
