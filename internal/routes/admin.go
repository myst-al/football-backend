package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(
	r *gin.RouterGroup,
	user *handler.UserHandler,
	team *handler.TeamHandler,
	player *handler.PlayerHandler,
) {
	UserAdminRoutes(r, user)
	TeamAdminRoutes(r, team)
	PlayerAdminRoutes(r, player)
}
