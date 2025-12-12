package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func GoalViewerRoutes(r *gin.RouterGroup, h *handler.GoalHandler) {
	r.GET("/goals/match/:match_id", h.GetByMatch)
	r.GET("/goals/top-scorers", h.TopScorers)
}

func GoalStaffRoutes(r *gin.RouterGroup, h *handler.GoalHandler) {
	r.POST("/goals", h.AddGoal)
}
