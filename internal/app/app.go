package app

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/api/handlers"
	"av-merch-shop/internal/repository"
	"av-merch-shop/internal/usecase"
	"av-merch-shop/pkg/auth"
	"av-merch-shop/pkg/database"
)

type Application struct {
	Handlers *handlers.Handlers
}

func New(cfg *config.Config, db *database.Database) *Application {
	repos := usecase.Repos{
		TransactionManager: db.TransactionManager,
		UserRepo:           repository.NewPGUserRepo(db, cfg.Logger),
		TransactionRepo:    repository.NewPGTransactionRepo(db, cfg.Logger),
		ItemRepo:           repository.NewPGItemRepo(db, cfg.Logger),
		OrderRepo:          repository.NewPGOrderRepo(db, cfg.Logger),
	}

	services := usecase.Services{
		TokenService: auth.NewJWTService(cfg),
		HashService:  auth.NewBCryptHashService(),
	}

	usecases := usecase.NewUsecases(repos, services)
	handlers := handlers.NewHandlers(cfg, usecases)

	return &Application{
		Handlers: handlers,
	}
}
