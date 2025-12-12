package routes

import (
	"football-backend/internal/handler"
	"football-backend/internal/middleware"
	"football-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterAll(
	r *gin.Engine,
	auth *handler.AuthHandler,
	user *handler.UserHandler,
	team *handler.TeamHandler,
	player *handler.PlayerHandler,
	match *handler.MatchHandler,
	goal *handler.GoalHandler,
	userRepo repository.UserRepository,
) {
	api := r.Group("/api/v1")

	AuthRoutes(api, auth, userRepo)

	secured := api.Group("/")
	secured.Use(middleware.JWTAuth(userRepo))

	viewer := secured.Group("/")
	viewer.Use(middleware.RequireRoles("ADMIN", "STAFF", "VIEWER"))
	RegisterViewerRoutes(viewer, user, team, player, match, goal)

	staff := secured.Group("/")
	staff.Use(middleware.RequireRoles("ADMIN", "STAFF"))
	RegisterStaffRoutes(staff, player, match, goal)

	admin := secured.Group("/")
	admin.Use(middleware.RequireRoles("ADMIN"))
	RegisterAdminRoutes(admin, user, team, player)
}
