package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	commoncfg "synapse/common/config"
)

var (
	Config *config
	DB     *gorm.DB
	Redis  redis.UniversalClient
)

type config struct {
	App        App                        `json:"app" yaml:"app" mapstructure:"app"`
	Logger     commoncfg.LoggerConfig     `json:"logger" yaml:"logger" mapstructure:"logger"`
	Server     commoncfg.ServerConfig     `json:"server" yaml:"server" mapstructure:"server"`
	Datasource commoncfg.DatasourceConfig `json:"datasource" yaml:"datasource" mapstructure:"datasource"`
	Redis      commoncfg.RedisConfig      `json:"redis" yaml:"redis" mapstructure:"redis"`
	GrpcServer commoncfg.GrpcServerConfig `json:"grpc_server" yaml:"grpc_server" mapstructure:"grpc_server"`
	Sentry     commoncfg.SentryConfig     `json:"sentry" yaml:"sentry" mapstructure:"sentry"`
}

type App struct {
	//AsyncApiWaitTimeout time.Duration             `json:"async_api_wait_timeout" yaml:"async_api_wait_timeout" mapstructure:"async_api_wait_timeout"`
	//AuthToken           string                    `json:"auth_token" yaml:"auth_token" mapstructure:"auth_token"`
	Services []commoncfg.ServiceConfig `json:"services" yaml:"services" mapstructure:"services"`
}
