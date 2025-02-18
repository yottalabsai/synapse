package routers

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"synapse/api/controllers"
	"synapse/api/middleware"
	"synapse/common"
	"synapse/config"
	"synapse/job"
	"synapse/log"
	"synapse/rpc"
	"synapse/service"
	"synapse/utils"
)

func InitRouter(ctx context.Context, engine *gin.Engine) error {
	engine.Use(ginzap.RecoveryWithZap(log.ZapLog, true))
	// init rpc client
	serviceConfigs := config.MustGetServiceConfig(common.ServiceYottaSaaS)
	yottaSaaSClient := rpc.NewYottaSaaSClient(&serviceConfigs[0], resty.NewWithClient(utils.ProxiedClientFromEnv()).SetLogger(log.Log))
	inferencePublicModelJob := job.NewInferencePublicModelJob(ctx, yottaSaaSClient)
	// Init other services
	svc := service.NewServerlessService(config.DB)
	// Health check
	// engine.GET("/actuator/health/liveness", healthHandler)
	var (
		// 分组
		apiGroupAuth = engine.Group("/api/v1", middleware.RequestHeader(), middleware.Authentication())
	)

	// 缓存处理
	//  tokenProvider = cache.NewDBTokenProvider(datasource.Db, datasource.WalletDB)
	//  poolProvider.Start(ctx, 30*time.Second)
	{
		ctl := controllers.NewServerlessController(svc)
		apiGroupAuth.GET("/endpoints/:endpointId", ctl.FindByEndpointId)
		apiGroupAuth.POST("/endpoints", ctl.CreateEndpoint)
		// 执行缓存清理
		// task.RunTransferTasks(ctx, svc)
	}

	{
		ctl := controllers.NewInferenceController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/inference", ctl.Inference)
		// 执行缓存清理
		// task.RunTransferTasks(ctx, svc)
	}

	{
		ctl := controllers.NewTextToImageController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/textToImage", ctl.Render)
		// 执行缓存清理
		// task.RunTransferTasks(ctx, svc)
	}

	{
		ctl := controllers.NewImageController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/images", ctl.Render)
		// 执行缓存清理
		// task.RunTransferTasks(ctx, svc)
	}

	// 启动定时任务
	jobManager := job.NewSynapseJobManager(inferencePublicModelJob)
	jobManager.StartJobs()

	return nil
}
