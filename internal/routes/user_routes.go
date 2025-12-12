package routes

import (
	"football-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func UserViewerRoutes(r *gin.RouterGroup, h *handler.UserHandler) {
	r.GET("/users/:id", h.GetByID)
}

func UserAdminRoutes(r *gin.RouterGroup, h *handler.UserHandler) {
	r.GET("/users", h.GetAdmins)
	r.DELETE("/users/:id", h.Delete)
}
