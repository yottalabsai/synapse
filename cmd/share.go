package cmd

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdlog "log"
	"os"
	"synapse/config"
	"synapse/log"
	"time"
)

var (
	logger *zap.Logger
)

var (
	Name    = "synapse"
	Version = "v2"
	Commit  = "N.A."
)

func PrintVersion() {
	log.Log.Infow("version info", "name", Name, "version", Version, "commit", Commit)
}

func InitLogger() {
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
