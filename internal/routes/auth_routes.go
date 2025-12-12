package routes

import (
	"football-backend/internal/handler"
	"football-backend/internal/middleware"
	"football-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h *handler.AuthHandler, userRepo repository.UserRepository) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/logout", h.Logout)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(userRepo))
	protected.GET("/me", h.Me)
}
