package main

import (
	_ "tz_ozon/docs"
	"tz_ozon/internal/config"
	"tz_ozon/internal/db"
	"tz_ozon/internal/logger"
	"tz_ozon/internal/router"

	"go.uber.org/zap"
)

// @Title TZ_OZON API
// @Version 1.0
func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}

	log := logger.L()
	log.Info("Инициализация логгера успешна")
	cfg := config.Load(log, ".env")

	db.ConnectDB(cfg, log)
	db.Migrate(log)

	repo := db.NewGormExchangeRateRepo(db.DB)
	r := router.Router(log, repo)
	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to run server", zap.Error(err))
	}
}
