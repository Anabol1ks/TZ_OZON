package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExchangeRateMock struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date      time.Time `gorm:"not null;index"`
	XMLBody   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *ExchangeRateMock) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
