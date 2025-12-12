package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func MatchViewerRoutes(r *gin.RouterGroup, h *handler.MatchHandler) {
	r.GET("/matches", h.GetAll)
	r.GET("/matches/:id", h.GetByID)
	r.GET("/matches/:id/report", h.Report)
	r.GET("/matches/standing", h.Standing)
}

func MatchStaffRoutes(r *gin.RouterGroup, h *handler.MatchHandler) {
	r.POST("/matches", h.Create)
	r.PUT("/matches/:id", h.Update)
	r.POST("/matches/:id/result", h.SubmitResult)
}
