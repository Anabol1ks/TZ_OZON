package db

import (
	"os"
	"tz_ozon/internal/models"

	"go.uber.org/zap"
)

func Migrate(log *zap.Logger) {
	if err := DB.AutoMigrate(&models.ExchangeRateMock{}); err != nil {
		log.Error("Ошибка при миграции таблиц", zap.Error(err))
		os.Exit(1)
	}
	log.Info("Автомиграция таблиц завершена успешно")

}
