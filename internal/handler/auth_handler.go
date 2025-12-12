package handler

import (
	"football-backend/internal/dto"
	"football-backend/internal/response"
	"football-backend/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	if err := h.service.Register(input.Username, input.Password, input.Role); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 201, "User berhasil dibuat", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "Input tidak valid")
		return
	}

	access, refresh, user, err := h.service.Login(input.Username, input.Password)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Login berhasil", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          dto.ToUserDTO(user),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uid.(uint)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Profil user berhasil diambil", dto.ToUserDTO(user))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, 400, "refresh_token missing")
		return
	}

	if strings.TrimSpace(body.RefreshToken) == "" {
		response.Error(c, 400, "refresh_token missing")
		return
	}

	access, newRefresh, user, err := h.service.Refresh(body.RefreshToken)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Token diperbarui", gin.H{
		"access_token":  access,
		"refresh_token": newRefresh,
		"user":          dto.ToUserDTO(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, 400, "refresh_token missing")
		return
	}

	if err := h.service.Logout(body.RefreshToken); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, 200, "Berhasil logout", nil)
}
