package dto

import "football-backend/internal/models"

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func ToUserDTO(u *models.User) UserDTO {
	return UserDTO{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}
