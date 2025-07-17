package router

import (
	"tz_ozon/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Router(log *zap.Logger, repo handler.ExchangeRateRepository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	h := handler.NewExchangeHandler(log, repo)
	r.GET("/scripts/XML_daily.asp", h.GetExchangeRate)

	return r
}
