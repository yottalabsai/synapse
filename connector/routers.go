package connector

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"synapse/common/log"
	"synapse/connector/controllers"
	"synapse/connector/middleware"
)

func InitRouter(ctx context.Context, engine *gin.Engine) error {
	engine.Use(ginzap.RecoveryWithZap(log.ZapLog, true))

	var (
		apiGroupAuth = engine.Group("/api/v1", middleware.RequestHeader(), middleware.Authentication())
	)

	{
		// Health check
		ctl := controllers.NewHealthController(nil)
		apiGroupAuth.GET("/health", ctl.Health)
		apiGroupAuth.GET("/status", ctl.Status)
	}

	return nil
}
