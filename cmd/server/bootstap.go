package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kittiphop/zombie_board_game_back/infrastructure/config"
	"github.com/kittiphop/zombie_board_game_back/internal/app/handler/http"
	"github.com/kittiphop/zombie_board_game_back/internal/app/usecase"
	"github.com/kittiphop/zombie_board_game_back/internal/platform/database"
	"gorm.io/gorm"
)

func initializeApp(cfg *config.Config) (*gin.Engine, error) {
	// ===== Setup Database (Postgres) =====
	db, err := setupDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// ===== Initialize Repositories & UseCases =====
	userRepo := database.NewUserPostGres(db)
	userUC := usecase.NewUserUseCase(userRepo)

	roomRepo := database.NewRoomPostGres(db)
	roomUseCase := usecase.NewRoomUseCase(roomRepo)

	// ===== Initialize Handlers =====
	http.InitializeHandlers(&http.HandlerDeps{
		UserUC: userUC,
		RoomUC: roomUseCase,
	})

	// ===== Setup Router =====
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	// ===== Middleware =====
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	router.Use(gin.Recovery())

	http.SetupRoutes(router)

	return router, nil
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
	// ===== Initialize Postgres =====
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port)

	// ===== Connect to Postgres =====
	db, err := database.InitializePostgres(dsn)
	if err != nil {
		return nil, err
	}

	// ===== Run Migrations =====
	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
