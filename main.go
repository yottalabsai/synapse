package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/reflection"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"synapse/config"
	"synapse/log"
	"synapse/repository/types"
	"synapse/routers"
	"syscall"
	"time"

	"synapse/service"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	ecszap "go.elastic.co/ecszap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger *zap.Logger
)

func initLogger() {
	time.Local = time.UTC
	var err error

	if !config.Config.Logger.Elasticsearch.Enabled {
		logger, err = zap.NewDevelopment()
		if err != nil {
			stdlog.Fatalf("can't initialize zap logger: %v", err)
		}
	} else {
		encoderConfig := ecszap.NewDefaultEncoderConfig()
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		developmentConfig := zap.NewDevelopmentConfig()
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), developmentConfig.Level)

		// Elasticsearch core
		hook, err := newElasticHook("synapse")
		if err != nil {
			stdlog.Fatalf("Failed to create elastic hook: %v", err)
		}
		elasticCore := ecszap.NewCore(encoderConfig, zapcore.AddSync(hook), zap.DebugLevel)
		// Combine cores
		core := zapcore.NewTee(consoleCore, elasticCore)
		logger = zap.New(core, zap.AddCaller())

	}

	slogger := logger.Sugar()
	log.Log = slogger

	log.ZapLog = logger
}

func main() {
	flag.Parse()

	defer logger.Sync()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("unknown error occurred", zap.Any("err", err))
		}
	}()

	_, err := config.ReadConfig(&config.Config)
	if err != nil {
		logger.Fatal("read config failed", zap.Error(err))
	}

	initLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	config.Redis = config.InitRedis(&config.Config.Redis)

	if err := config.InitDatasource(ctx, logger, &config.Config.Datasource); err != nil {
		logger.Fatal("init datasource failed", zap.Error(err))
	}

	MigrateDB()

	if err := Start(ctx); err != nil {
		logger.Fatal("start app failed", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("app service stopped")
}

func MigrateDB() {
	err := config.DB.AutoMigrate(&types.ServerlessResource{})
	if err != nil {
		logger.Error("db migration error", zap.Error(err))
	}
}

func Start(ctx context.Context) error {

	engine := gin.New()
	if config.Config.Server.GIN.Mode != "" {
		gin.SetMode(config.Config.Server.GIN.Mode)
	}
	engine.Use(gin.Recovery())

	if err := routers.InitRouter(ctx, engine); err != nil {
		return errors.WithMessagef(err, "init router error")
	}

	go func() {
		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port),
			engine,
		); err != nil {
			log.Log.Warnw("http server stopped", zap.Error(err))
		}
	}()

	// start grpc server
	go func() {
		if err := startGrpc(); err != nil {
			logger.Error("grpc server stopped", zap.Error(err))
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

	log.Log.Infow("grpc server listening at:", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Log.Errorw("failed to serve:", err)
		return err
	}
	return nil
}

type elasticHook struct {
	client *elasticsearch.Client
	index  string
}

func (h *elasticHook) Write(p []byte) (n int, err error) {
	_, err = h.client.Index(
		h.index,
		bytes.NewReader(p),
		h.client.Index.WithContext(context.Background()),
		h.client.Index.WithDocumentType("_doc"),
		h.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func newElasticHook(index string) (*elasticHook, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Config.Logger.Elasticsearch.Host,
		},
		Username: config.Config.Logger.Elasticsearch.Username,
		Password: config.Config.Logger.Elasticsearch.Password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &elasticHook{client: es, index: index}, nil
}
