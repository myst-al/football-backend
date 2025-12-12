package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func PlayerViewerRoutes(r *gin.RouterGroup, h *handler.PlayerHandler) {
	r.GET("/players", h.GetAll)
	r.GET("/players/:id", h.GetByID)
	r.GET("/players/by-team/:team_id", h.GetByTeam)
}

func PlayerStaffRoutes(r *gin.RouterGroup, h *handler.PlayerHandler) {
	r.POST("/players", h.Create)
	r.PUT("/players/:id", h.Update)
	r.POST("/players/:id/transfer", h.Transfer)
}

func PlayerAdminRoutes(r *gin.RouterGroup, h *handler.PlayerHandler) {
	r.DELETE("/players/:id", h.Delete)
}
