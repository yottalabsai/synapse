package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os/signal"
	"synapse/common"
	config2 "synapse/common/config"
	"synapse/connector/service"
	"synapse/log"
	"synapse/worker"
	"synapse/worker/repository/types"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"synapse/cmd"
)

func main() {
	flag.Parse()

	defer log.ZapLog.Sync()
	defer func() {
		if err := recover(); err != nil {
			log.ZapLog.Error("unknown error occurred", zap.Any("err", err))
		}
	}()

	_, err := config2.ReadConfig(common.ServiceConnector, &config2.Config)
	if err != nil {
		log.ZapLog.Fatal("read config failed", zap.Error(err))
	}

	cmd.InitLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	config2.Redis = config2.InitRedis(&config2.Config.Redis)

	if err := config2.InitDatasource(ctx, log.ZapLog, &config2.Config.Datasource); err != nil {
		log.ZapLog.Fatal("init datasource failed", zap.Error(err))
	}

	MigrateDB()

	if err := Start(ctx); err != nil {
		log.ZapLog.Fatal("start app failed", zap.Error(err))
	}

	<-ctx.Done()
	log.ZapLog.Info("app service stopped")
}

func MigrateDB() {
	err := config2.DB.AutoMigrate(&types.ServerlessResource{})
	if err != nil {
		log.ZapLog.Error("db migration error", zap.Error(err))
	}
}

func Start(ctx context.Context) error {

	engine := gin.New()
	if config2.Config.Server.GIN.Mode != "" {
		gin.SetMode(config2.Config.Server.GIN.Mode)
	}
	engine.Use(gin.Recovery())

	if err := worker.InitRouter(ctx, engine); err != nil {
		return errors.WithMessagef(err, "init router error")
	}

	go func() {
		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", config2.Config.Server.Host, config2.Config.Server.Port),
			engine,
		); err != nil {
			log.Log.Warnw("http server stopped", zap.Error(err))
		}
	}()

	// start grpc server
	go func() {
		if err := startGrpc(); err != nil {
			log.ZapLog.Error("grpc server stopped", zap.Error(err))
		}
	}()

	return nil
}

func startGrpc() error {
	grpcServerConfig := config2.Config.GrpcServer
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcServerConfig.Host, grpcServerConfig.Port))
	if err != nil {
		log.Log.Errorw("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	// Register reflection service on gRPC server.
	reflection.Register(s)
	synapseGrpc.RegisterSynapseServiceServer(s, service.NewSynapseServer())

	log.Log.Infow("grpc server listening at:", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Log.Errorw("failed to serve:", err)
		return err
	}
	return nil
}
