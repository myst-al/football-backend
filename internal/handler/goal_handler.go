package handler

import (
	"football-backend/internal/dto"
	"football-backend/internal/models"
	"football-backend/internal/response"
	"football-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GoalHandler struct {
	service service.GoalService
}

func NewGoalHandler(s service.GoalService) *GoalHandler {
	return &GoalHandler{s}
}

func (h *GoalHandler) AddGoal(c *gin.Context) {
	var input struct {
		MatchID        uint   `json:"match_id" binding:"required"`
		TeamID         uint   `json:"team_id" binding:"required"`
		ScorerPlayerID uint   `json:"scorer_player_id" binding:"required"`
		Minute         string `json:"minute" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	goal := models.Goal{
		MatchID:        input.MatchID,
		TeamID:         input.TeamID,
		ScorerPlayerID: input.ScorerPlayerID,
		Minute:         input.Minute,
	}

	if err := h.service.AddGoal(&goal); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 201, "Gol berhasil ditambahkan", nil)
}

func (h *GoalHandler) GetByMatch(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("match_id"))

	list, err := h.service.GetGoals(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	dtoList := []dto.GoalDTO{}
	for _, g := range list {
		dtoList = append(dtoList, dto.ToGoalDTO(&g))
	}

	response.Success(c, 200, "Data gol berhasil diambil", dtoList)
}

func (h *GoalHandler) TopScorers(c *gin.Context) {
	list, err := h.service.TopScorers(10)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Top skor berhasil diambil", list)
}
