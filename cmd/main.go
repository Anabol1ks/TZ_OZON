package main

import (
	"tz_ozon/internal/config"
	"tz_ozon/internal/db"
	"tz_ozon/internal/logger"
	"tz_ozon/internal/router"

	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}

	log := logger.L()
	log.Info("Инициализация логгера успешна")
	cfg := config.Load(log)

	db.ConnectDB(cfg, log)
	db.Migrate(log)

	r := router.Router(cfg, log)
	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to run server", zap.Error(err))
	}
}
