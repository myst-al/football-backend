package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterViewerRoutes(
	r *gin.RouterGroup,
	user *handler.UserHandler,
	team *handler.TeamHandler,
	player *handler.PlayerHandler,
	match *handler.MatchHandler,
	goal *handler.GoalHandler,
) {
	UserViewerRoutes(r, user)
	TeamViewerRoutes(r, team)
	PlayerViewerRoutes(r, player)
	MatchViewerRoutes(r, match)
	GoalViewerRoutes(r, goal)
}
