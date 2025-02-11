package app

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/handlers"
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
		UserRepo: repository.NewPGUserRepo(db),
	}

	services := usecase.Services{
		TokenService: auth.NewJWTService(cfg),
	}

	usecases := usecase.NewUsecases(repos, services)
	handlers := handlers.NewHandlers(cfg, usecases)

	return &Application{
		Handlers: handlers,
	}
}
