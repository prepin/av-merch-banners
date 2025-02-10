package app

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/handlers"
)

type Application struct {
	Handlers *handlers.Handlers
}

func New(cfg *config.Config) *Application {
	handlers := handlers.NewHandlers(cfg)

	return &Application{
		Handlers: handlers,
	}
}
