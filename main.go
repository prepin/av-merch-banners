package main

import (
	"av-merch-shop/config"
)

func main() {
	cfg := config.Load()

	cfg.Logger.Info("Merch Shop started")
}
