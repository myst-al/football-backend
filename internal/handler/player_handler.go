package handler

import (
	"football-backend/internal/dto"
	"football-backend/internal/models"
	"football-backend/internal/response"
	"football-backend/internal/service"
	"football-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService service.PlayerService
}

func NewPlayerHandler(s service.PlayerService) *PlayerHandler {
	return &PlayerHandler{s}
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var input struct {
		TeamID       uint   `json:"team_id" binding:"required"`
		Name         string `json:"name" binding:"required"`
		HeightCM     int    `json:"height"`
		WeightKG     int    `json:"weight"`
		Position     string `json:"position" binding:"required"`
		JerseyNumber int    `json:"jersey_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	p := models.Player{
		TeamID:       input.TeamID,
		Name:         input.Name,
		HeightCM:     input.HeightCM,
		WeightKG:     input.WeightKG,
		Position:     input.Position,
		JerseyNumber: input.JerseyNumber,
	}

	if err := h.playerService.Create(&p); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 201, "Pemain berhasil dibuat", dto.ToPlayerDTO(&p))
}

func (h *PlayerHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	p, err := h.playerService.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	if err := c.ShouldBindJSON(p); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	if err := h.playerService.Update(p); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Pemain berhasil diperbarui", dto.ToPlayerDTO(p))
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.playerService.Delete(uint(id)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Pemain berhasil dihapus", nil)
}

func (h *PlayerHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	p, err := h.playerService.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Data pemain berhasil diambil", dto.ToPlayerDTO(p))
}

func (h *PlayerHandler) GetByTeam(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("team_id"))

	list, err := h.playerService.GetByTeam(uint(teamID))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Data pemain berhasil diambil", dto.ToPlayerDTOList(list))
}

func (h *PlayerHandler) Transfer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		NewTeamID    uint `json:"new_team_id" binding:"required"`
		JerseyNumber int  `json:"jersey_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	if err := h.playerService.TransferPlayer(uint(id), input.NewTeamID, input.JerseyNumber); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Transfer pemain berhasil", nil)
}

func (h *PlayerHandler) GetAll(c *gin.Context) {
	q := utils.ParseQuery(c)

	result, err := h.playerService.GetList(q)
	if err != nil {
		response.FromError(c, err)
		return
	}

	items := result["items"].([]models.Player)

	dtoList := dto.ToPlayerDTOList(items)
	result["items"] = dtoList

	response.Success(c, 200, "Data pemain berhasil diambil", result)
}
