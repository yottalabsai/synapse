package connector

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"synapse/common"
	commonConfig "synapse/common/config"
	"synapse/common/utils"
	"synapse/connector/config"
	"synapse/connector/controllers"
	"synapse/connector/middleware"
	"synapse/connector/rpc"
	"synapse/worker/job"
)

func InitRouter(ctx context.Context, engine *gin.Engine) error {
	// init rpc client
	serviceConfigs := commonConfig.MustGetServiceConfig(config.Config.App.Services, common.ServiceYottaSaaS)
	yottaSaaSClient := rpc.NewYottaSaaSClient(&serviceConfigs[0], resty.NewWithClient(utils.ProxiedClientFromEnv()))
	inferencePublicModelJob := job.NewInferencePublicModelJob(ctx, yottaSaaSClient)
	// Init other services

	var (
		apiGroupAuth = engine.Group("/api/v1", middleware.RequestHeader(), middleware.Authentication())
	)

	{
		// Health check
		ctl := controllers.NewHealthController(nil)
		apiGroupAuth.GET("/health", ctl.Health)
		apiGroupAuth.GET("/status", ctl.Status)
	}

	// Run job
	jobManager := job.NewSynapseJobManager(inferencePublicModelJob)
	jobManager.StartJobs()

	return nil
}
