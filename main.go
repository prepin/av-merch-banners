package main

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/app"
	"av-merch-shop/pkg/server"
	"log/slog"
)

func main() {
	cfg := config.Load()
	app := app.New(cfg)

	cfg.Logger.Info("Launching Merch Shop", "config", cfg.Server)

	srv := server.New(cfg, app.Handlers)
	if err := srv.Run(); err != nil {
		slog.Error("failed to launch server", "Error", err)
	}
}
