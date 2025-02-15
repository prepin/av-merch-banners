package main

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/app"
	"av-merch-shop/pkg/database"
	"av-merch-shop/pkg/redis"
	"av-merch-shop/pkg/server"
	"log/slog"
)

func main() {
	cfg := config.Load()
	db := database.NewDatabase(cfg.DB)
	defer db.Close()

	redis := redis.NewRedis(cfg.Redis, cfg.Logger)

	app := app.New(cfg, db, redis)

	cfg.Logger.Info("Launching Merch Shop", "config", cfg.Server)

	srv := server.New(cfg, app.Handlers)
	if err := srv.Run(); err != nil {
		slog.Error("failed to launch server", "Error", err)
	}
}
