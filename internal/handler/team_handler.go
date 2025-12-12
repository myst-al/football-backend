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

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(s service.TeamService) *TeamHandler {
	return &TeamHandler{s}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		LogoURL     string `json:"logo_url"`
		YearFounded int    `json:"year_founded"`
		Address     string `json:"address"`
		City        string `json:"city"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	team := models.Team{
		Name:        input.Name,
		LogoURL:     input.LogoURL,
		YearFounded: input.YearFounded,
		Address:     input.Address,
		City:        input.City,
	}

	if err := h.service.Create(&team); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 201, "Team berhasil dibuat", dto.ToTeamDTO(&team))
}

func (h *TeamHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	t, err := h.service.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	if err := c.ShouldBindJSON(t); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	if err := h.service.Update(t); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Team berhasil diperbarui", dto.ToTeamDTO(t))
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Team berhasil dihapus", nil)
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	t, err := h.service.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Data team berhasil diambil", dto.ToTeamDTO(t))
}

func (h *TeamHandler) GetAll(c *gin.Context) {
	q := utils.ParseQuery(c)

	result, err := h.service.GetList(q)
	if err != nil {
		response.FromError(c, err)
		return
	}

	items := result["items"].([]models.Team)
	dtoList := dto.ToTeamDTOList(items)
	result["items"] = dtoList

	response.Success(c, 200, "Data tim berhasil diambil", result)
}
