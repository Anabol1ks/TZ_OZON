package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tz_ozon/internal/handler"
	"tz_ozon/internal/logger"
	"tz_ozon/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockRepo struct {
	findByDateFunc func(date time.Time) (*models.ExchangeRateMock, error)
	createFunc     func(rate *models.ExchangeRateMock) error
}

func (m *mockRepo) FindByDate(date time.Time) (*models.ExchangeRateMock, error) {
	return m.findByDateFunc(date)
}
func (m *mockRepo) Create(rate *models.ExchangeRateMock) error {
	return m.createFunc(rate)
}

func TestGetExchangeRate_Success(t *testing.T) {
	_ = logger.Init()
	log := logger.L()
	log.Info("Тест: успешное получение курса валют")
	repo := &mockRepo{
		findByDateFunc: func(date time.Time) (*models.ExchangeRateMock, error) {
			log.Info("mockRepo.FindByDate вызван", zap.Time("date", date))
			return &models.ExchangeRateMock{
				Date:    date,
				XMLBody: "<ValCurs><Valute>Foreign Currency Market</Valute></ValCurs>",
			}, nil
		},
		createFunc: func(rate *models.ExchangeRateMock) error {
			log.Info("mockRepo.Create вызван", zap.Time("date", rate.Date))
			return nil
		},
	}
	h := handler.NewExchangeHandler(log, repo)
	r := ginTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/scripts/XML_daily.asp?date_req=17/07/2025", nil)
	log.Info("Отправка GET-запроса на /scripts/XML_daily.asp?date_req=17/07/2025")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Error("Ожидался статус 200", zap.Int("code", w.Code), zap.String("body", w.Body.String()))
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<ValCurs")
	assert.Contains(t, w.Body.String(), "Valute")
	assert.Contains(t, w.Body.String(), "Foreign Currency Market")
}

func TestGetExchangeRate_BadRequest(t *testing.T) {
	_ = logger.Init()
	log := logger.L()
	log.Info("Тест: отсутствие параметра date_req")
	repo := &mockRepo{}
	h := handler.NewExchangeHandler(log, repo)
	r := ginTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/scripts/XML_daily.asp", nil)
	log.Info("Отправка GET-запроса без date_req")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		log.Error("Ожидался статус 400", zap.Int("code", w.Code), zap.String("body", w.Body.String()))
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "date_req is required")
}

func TestGetExchangeRate_InvalidDate(t *testing.T) {
	_ = logger.Init()
	log := logger.L()
	log.Info("Тест: некорректный формат даты")
	repo := &mockRepo{}
	h := handler.NewExchangeHandler(log, repo)
	r := ginTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/scripts/XML_daily.asp?date_req=2025-07-17", nil)
	log.Info("Отправка GET-запроса с некорректной датой")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		log.Error("Ожидался статус 400", zap.Int("code", w.Code), zap.String("body", w.Body.String()))
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid date format")
}

// gin router для теста
func ginTestRouter(h *handler.ExchangeHandler) *gin.Engine {
	r := gin.New()
	r.GET("/scripts/XML_daily.asp", h.GetExchangeRate)
	return r
}
