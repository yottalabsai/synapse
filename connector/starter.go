package connector

import (
	"context"
	"fmt"
	sentryGin "github.com/cockroachdb/sentry-go/gin"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	idgenerator "synapse/common/id-generator"
	"synapse/common/log"
	"synapse/connector/config"
	"synapse/connector/service"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Start(ctx context.Context) error {

	if err := idgenerator.InitMultiSnowflakeInstances(
		ctx,
		config.Redis,
		config.SnowflakeNodeIDRedisKeyForNextId,
	); err != nil {
		return errors.WithMessagef(err, "init multi snowflake instances failed")
	}

	engine := gin.New()
	if config.Config.Server.GIN.Mode != "" {
		gin.SetMode(config.Config.Server.GIN.Mode)
	}
	engine.Use(gin.Recovery())

	if config.Config.Sentry.Enabled {
		engine.Use(sentryGin.New(sentryGin.Options{})) // use sentry
	}

	if err := InitRouter(ctx, engine); err != nil {
		return errors.WithMessagef(err, "init router failed")
	}

	go func() {
		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port),
			engine,
		); err != nil {
			log.Log.Warn("http server stopped", zap.Error(err))
		}
	}()

	// start grpc server
	go func() {
		if err := startGrpc(); err != nil {
			log.Log.Error("grpc server stopped", zap.Error(err))
		}
	}()

	return nil
}

func startGrpc() error {
	grpcServerConfig := config.Config.GrpcServer
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcServerConfig.Host, grpcServerConfig.Port))
	if err != nil {
		log.Log.Errorw("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	// Register reflection service on gRPC server.
	reflection.Register(s)
	synapseGrpc.RegisterSynapseServiceServer(s, service.NewSynapseServer())

	log.Log.Infof("grpc server listening at: %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Log.Errorw("failed to serve:", err)
		return err
	}
	return nil
}
