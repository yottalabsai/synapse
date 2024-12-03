package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	stdlog "log"
	"net/http"
	"os/signal"
	"synapse/config"
	"synapse/log"
	"synapse/routers"
	"syscall"
	"time"
)

var (
	logger *zap.Logger
)

func init() {
	time.Local = time.UTC
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		stdlog.Fatalf("can't initialize zap logger: %v", err)
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	config.Redis = config.InitRedis(&config.Config.Redis)

	if err := config.InitDatasource(ctx, logger, &config.Config.Datasource); err != nil {
		logger.Fatal("init datasource failed", zap.Error(err))
	}

	if err := Start(ctx); err != nil {
		logger.Fatal("start app failed", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("app service stopped")
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
			log.Log.Warn("http server stopped", zap.Error(err))
		}
	}()

	return nil
}
