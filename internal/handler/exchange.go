package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"tz_ozon/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ExchangeRateRepository interface {
	FindByDate(date time.Time) (*models.ExchangeRateMock, error)
	Create(rate *models.ExchangeRateMock) error
}

type ExchangeHandler struct {
	Log  *zap.Logger
	Repo ExchangeRateRepository
}

func NewExchangeHandler(log *zap.Logger, repo ExchangeRateRepository) *ExchangeHandler {
	return &ExchangeHandler{
		Log:  log,
		Repo: repo,
	}
}

// GetExchangeRate godoc
// @Summary      Получить курс валют
// @Description  Мокаем ответ от ЦБ РФ на конкретную дату
// @Tags         exchange
// @Param        date_req  query  string  true  "Дата в формате dd/mm/yyyy"  example(02/03/2023)
// @Produce      xml
// @Success      200  {object}  string  "Успешный XML ответ"
// @Failure      400  {string}  string  "Некорректный запрос"
// @Failure      500  {string}  string  "Внутренняя ошибка"
// @Router       /scripts/XML_daily.asp [get]
func (h *ExchangeHandler) GetExchangeRate(c *gin.Context) {
	dateStr := c.Query("date_req")
	if dateStr == "" {
		c.String(http.StatusBadRequest, "date_req is required")
		return
	}

	date, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid date format")
		return
	}

	rate, err := h.Repo.FindByDate(date)

	switch {
	case err == nil && rate != nil:
		c.Data(http.StatusOK, "application/json", []byte(rate.XMLBody))
		return

	case err == gorm.ErrRecordNotFound || rate == nil:
		if rand.Intn(2) == 0 {
			h.Log.Warn("Симуляция ошибки 500", zap.String("date_req", dateStr))
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		xml := generateFakeXML(date)
		newRate := &models.ExchangeRateMock{
			Date:    date,
			XMLBody: xml,
		}

		if err := h.Repo.Create(newRate); err != nil {
			h.Log.Error("Не удалось сохранить сгенерированный XML-файл", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to persist data")
			return
		}

		c.Data(http.StatusOK, "application/json", []byte(newRate.XMLBody))
		return

	default:
		h.Log.Error("Неизвестная ошибка", zap.Error(err))
		c.String(http.StatusInternalServerError, "DB Error")
		return
	}
}

func generateFakeXML(date time.Time) string {
	type currency struct {
		ID       string
		NumCode  string
		CharCode string
		Nominal  int
		Name     string
	}
	currencies := []currency{
		{"R01010", "036", "AUD", 1, "Австралийский доллар"},
		{"R01020A", "944", "AZN", 1, "Азербайджанский манат"},
		{"R01035", "826", "GBP", 1, "Фунт стерлингов"},
	}

	xml := "<?xml version=\"1.0\" encoding=\"windows-1251\"?>\n"
	xml += "<ValCurs Date=\"" + date.Format("02.01.2006") + "\" name=\"Foreign Currency Market\">\n"
	for _, c := range currencies {
		value := rand.Float64()*60 + 30 // 30..90
		rate := formatRateComma(value)
		xml += "  <Valute ID=\"" + c.ID + "\">\n"
		xml += "    <NumCode>" + c.NumCode + "</NumCode>\n"
		xml += "    <CharCode>" + c.CharCode + "</CharCode>\n"
		xml += "    <Nominal>" + fmt.Sprintf("%d", c.Nominal) + "</Nominal>\n"
		xml += "    <Name>" + c.Name + "</Name>\n"
		xml += "    <Value>" + rate + "</Value>\n"
		xml += "    <VunitRate>" + rate + "</VunitRate>\n"
		xml += "  </Valute>\n"
	}
	xml += "</ValCurs>"
	return xml
}

func formatRateComma(value float64) string {
	s := fmt.Sprintf("%.4f", value)
	return replaceDotWithComma(s)
}

func replaceDotWithComma(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return s[:i] + "," + s[i+1:]
		}
	}
	return s
}
