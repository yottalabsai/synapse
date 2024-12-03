package routers

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"synapse/api/controllers"
	"synapse/api/middleware"
	"synapse/config"
	"synapse/log"
	"synapse/service"
)

func InitRouter(ctx context.Context, engine *gin.Engine) error {
	engine.Use(ginzap.RecoveryWithZap(log.ZapLog, true))
	// Init other services
	svc := service.NewSwapService(config.DB)
	// Health check
	// engine.GET("/actuator/health/liveness", healthHandler)
	var (
		// 分组
		apiGroupAuth = engine.Group("/api", middleware.RequestHeader(), middleware.Authentication())
	)

	// 缓存处理
	//  tokenProvider = cache.NewDBTokenProvider(datasource.Db, datasource.WalletDB)
	//  poolProvider.Start(ctx, 30*time.Second)
	{
		ctl := controllers.NewDemo(svc)
		apiGroupAuth.GET("/demo", ctl.Demo)
		// 执行缓存清理
		// task.RunTransferTasks(ctx, svc)
	}

	return nil
}
