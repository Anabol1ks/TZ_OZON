package db

import (
	"time"
	"tz_ozon/internal/models"

	"gorm.io/gorm"
)

type GormExchangeRateRepo struct {
	DB *gorm.DB
}

func NewGormExchangeRateRepo(db *gorm.DB) *GormExchangeRateRepo {
	return &GormExchangeRateRepo{DB: db}
}

func (r *GormExchangeRateRepo) FindByDate(date time.Time) (*models.ExchangeRateMock, error) {
	var rate models.ExchangeRateMock
	err := r.DB.Where("date = ?", date).First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *GormExchangeRateRepo) Create(rate *models.ExchangeRateMock) error {
	return r.DB.Create(rate).Error
}
