package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func TeamViewerRoutes(r *gin.RouterGroup, h *handler.TeamHandler) {
	r.GET("/teams", h.GetAll)
	r.GET("/teams/:id", h.GetByID)
}

func TeamAdminRoutes(r *gin.RouterGroup, h *handler.TeamHandler) {
	r.POST("/teams", h.Create)
	r.PUT("/teams/:id", h.Update)
	r.DELETE("/teams/:id", h.Delete)
}
