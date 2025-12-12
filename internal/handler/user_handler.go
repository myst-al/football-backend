package handler

import (
	"football-backend/internal/response"
	"football-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{us}
}

func (h *UserHandler) GetAdmins(c *gin.Context) {
	list, err := h.userService.GetAdmins()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, 200, "Data admin berhasil diambil", list)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	admin, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Data admin berhasil diambil", admin)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.userService.Delete(uint(id)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Admin berhasil dihapus", nil)
}
