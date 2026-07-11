package main

import (
	"log"

	"wrappedweekly/backend/internal/aiprovider"
	"wrappedweekly/backend/internal/config"
	"wrappedweekly/backend/internal/handler"
	"wrappedweekly/backend/internal/repository"
	"wrappedweekly/backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}

	pool, err := repository.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek ke database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	activityRepo := repository.NewActivityRepository(pool)
	recapRepo := repository.NewRecapRepository(pool)

	jwtManager := usecase.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtManager)
	activityUsecase := usecase.NewActivityUsecase(activityRepo)

	aiProvider := aiprovider.NewProvider(cfg.AIProvider)
	recapUsecase := usecase.NewRecapUsecase(recapRepo, activityRepo, userRepo, aiProvider)

	handlers := handler.Handlers{
		Auth:      handler.NewAuthHandler(authUsecase, cfg),
		Activity:  handler.NewActivityHandler(activityUsecase),
		Recap:     handler.NewRecapHandler(recapUsecase),
		Dashboard: handler.NewDashboardHandler(activityRepo),
	}

	router := handler.NewRouter(cfg, jwtManager, handlers)

	log.Printf("Wrapped Weekly API listening on :%s (AI_PROVIDER=%s)", cfg.Port, cfg.AIProvider)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}
