package main

import (
	"log"

	"football-backend/internal/config"
	"football-backend/internal/database"
	"football-backend/internal/handler"
	"football-backend/internal/middleware"
	"football-backend/internal/models"
	"football-backend/internal/repository"
	"football-backend/internal/routes"
	"football-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func mustLoadEnv() {
	_ = godotenv.Load()
}

func main() {
	mustLoadEnv()

	cfg := config.Load()

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("database connect error: %v", err)
	}

	db := database.DB
	if db == nil {
		log.Fatal("database.DB is nil after Connect()")
	}

	if err := db.AutoMigrate(
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.Goal{},
		&models.User{},
		&models.PlayerTransfer{},
		&models.RefreshToken{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	playerTransferRepo := repository.NewPlayerTransferRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)

	teamSvc := service.NewTeamService(teamRepo)
	playerSvc := service.NewPlayerService(playerRepo, playerTransferRepo, teamRepo)
	goalSvc := service.NewGoalService(goalRepo, matchRepo)
	matchSvc := service.NewMatchService(matchRepo, goalRepo, teamRepo)
	authSvc := service.NewAuthService(userRepo, refreshRepo)
	userSvc := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	playerHandler := handler.NewPlayerHandler(playerSvc)
	goalHandler := handler.NewGoalHandler(goalSvc)
	matchHandler := handler.NewMatchHandler(matchSvc, goalSvc)

	r := gin.New()
	r.Use(middleware.JSONLogger())
	r.Use(gin.Recovery())

	routes.RegisterAll(
		r,
		authHandler,
		userHandler,
		teamHandler,
		playerHandler,
		matchHandler,
		goalHandler,
		userRepo,
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := ":" + cfg.AppPort
	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
