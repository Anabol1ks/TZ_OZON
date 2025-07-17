package router

import (
	"tz_ozon/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router(cfg *config.Config, log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}
