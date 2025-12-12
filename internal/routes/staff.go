package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStaffRoutes(
	r *gin.RouterGroup,
	player *handler.PlayerHandler,
	match *handler.MatchHandler,
	goal *handler.GoalHandler,
) {
	PlayerStaffRoutes(r, player)
	MatchStaffRoutes(r, match)
	GoalStaffRoutes(r, goal)
}
