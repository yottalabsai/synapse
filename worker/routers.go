package worker

import (
	"context"
	"github.com/gin-gonic/gin"
	"synapse/common/config"
	"synapse/worker/controllers"
	"synapse/worker/middleware"
	service "synapse/worker/service"
)

func InitRouter(ctx context.Context, engine *gin.Engine) error {
	// init rpc client
	// Init other services
	svc := service.NewServerlessService(config.DB)

	var (
		apiGroupAuth = engine.Group("/api/v1", middleware.RequestHeader(), middleware.Authentication())
	)

	{
		// Health check
		ctl := controllers.NewHealthController(nil)
		apiGroupAuth.GET("/health", ctl.Health)
		apiGroupAuth.GET("/status", ctl.Status)
	}

	{
		ctl := controllers.NewServerlessController(svc)
		apiGroupAuth.GET("/endpoints/:endpointId", ctl.FindByEndpointId)
		apiGroupAuth.POST("/endpoints", ctl.CreateEndpoint)
		// clean cache
		// task.RunTransferTasks(ctx, svc)
	}

	{
		ctl := controllers.NewInferenceController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/inference", ctl.Inference)
	}

	{
		ctl := controllers.NewTextToImageController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/textToImage", ctl.Render)
	}

	{
		ctl := controllers.NewImageController(config.GrpcServer)
		apiGroupAuth.POST("/endpoints/:endpointId/images", ctl.Render)
	}

	return nil
}
