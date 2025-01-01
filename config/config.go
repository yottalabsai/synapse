package config

import (
	"os"
	"reflect"
	"strings"
	"synapse/log"
	"synapse/service"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Config     *config
	DB         *gorm.DB
	Redis      redis.UniversalClient
	GrpcServer *service.SynapseServer
)

type config struct {
	Logger                  LoggerConfig     `json:"logger" yaml:"logger" mapstructure:"logger"`
	Server                  ServerConfig     `json:"server" yaml:"server" mapstructure:"server"`
	Redis                   RedisConfig      `json:"redis" yaml:"redis" mapstructure:"redis"`
	Datasource              DatasourceConfig `json:"datasource" yaml:"datasource" mapstructure:"datasource"`
	ApiKey                  string           `json:"api_key" yaml:"api_key" mapstructure:"api_key"`
	LockExpirationDuration  time.Duration    `json:"lock_expiration_duration" yaml:"lock_expiration_duration" mapstructure:"lock_expiration_duration"`
	AcquireLockSpinDuration time.Duration    `json:"acquire_lock_spin_duration" yaml:"acquire_lock_spin_duration" mapstructure:"acquire_lock_spin_duration"`
	GrpcServer              GrpcServerConfig `json:"grpc_server" yaml:"grpc_server" mapstructure:"grpc_server"`
	Services                []ServiceConfig  `json:"services" yaml:"services" mapstructure:"services"`
}

type ServerConfig struct {
	Host string    `json:"host" yaml:"host" mapstructure:"host"` // HTTP listen host
	Port int       `json:"port" yaml:"port" mapstructure:"port"` // HTTP listen port
	GIN  GINConfig `json:"gin" yaml:"gin" mapstructure:"gin"`    // GIN config

	Pprof PprofConfig `json:"pprof" yaml:"pprof" mapstructure:"pprof"`
}

type GINConfig struct {
	Mode string `json:"mode" yaml:"mode"` // GIN mode
}

type PprofConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"` // pprof enabled
	Prefix  string `json:"prefix" yaml:"prefix" mapstructure:"prefix"`    // pprof URL prefix
}

type KafkaConfig struct {
	Brokers []string `json:"brokers" yaml:"brokers" mapstructure:"brokers"`
}

type DatasourceConfig struct {
	Pool     DatasourcePoolConfig `json:"pool,omitempty" yaml:"pool" mapstructure:"pool"`
	GORM     GORMConfig           `json:"gorm,omitempty" yaml:"gorm" mapstructure:"gorm"`
	Postgres PostgresConfig       `json:"postgres,omitempty" yaml:"postgres" mapstructure:"postgres"`
}

type DatasourcePoolConfig struct {
	MaxOpenConns    int           `json:"max_open_conns,omitempty" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns,omitempty" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime,omitempty" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}

type GORMConfig struct {
	LogLevel string `json:"log_level,omitempty" yaml:"log_level" mapstructure:"log_level"`
}

type PostgresConfig struct {
	Host       string `json:"host,omitempty" yaml:"host" mapstructure:"host"`
	Port       int    `json:"port,omitempty" yaml:"port" mapstructure:"port"`
	Username   string `json:"username,omitempty" yaml:"username" mapstructure:"username"`
	Password   string `json:"password,omitempty" yaml:"password" mapstructure:"password"`
	Database   string `json:"database,omitempty" yaml:"database" mapstructure:"database"`
	SearchPath string `json:"search_path,omitempty" yaml:"search_path" mapstructure:"search_path"`
	TimeZone   string `json:"time_zone,omitempty" yaml:"time_zone" mapstructure:"time_zone"`
	SslMode    string `json:"ssl_mode,omitempty" yaml:"ssl_mode" mapstructure:"ssl_mode"`
}

type LoggerConfig struct {
	ConsoleColorEnabled bool                   `json:"console_color_enabled" yaml:"console_color_enabled" mapstructure:"console_color_enabled"` // console color enabled
	Appenders           []LoggerAppenderConfig `json:"appenders" yaml:"appenders" mapstructure:"appenders"`                                     // log appenders, default is console
}

type LoggerAppenderConfig struct {
	Level  string                   `json:"level" yaml:"level" mapstructure:"level"`    // log level, default is info
	Type   string                   `json:"type" yaml:"type" mapstructure:"type"`       // log type, support console, file, default is console
	Format string                   `json:"format" yaml:"format" mapstructure:"format"` // log format, support text, json, default is text
	File   LoggerAppenderFileConfig `json:"file" yaml:"file" mapstructure:"file"`       // log file config
}

type LoggerAppenderFileConfig struct {
	Filename         string        `json:"filename" yaml:"filename" mapstructure:"filename"`                            // log file name, default is 'app.log'
	FilenameSuffix   string        `json:"filename_suffix" yaml:"filename_suffix" mapstructure:"filename_suffix"`       // log file name suffix, default is '.%Y%m%d'
	MaxAgeDuration   time.Duration `json:"max_age_duration" yaml:"max_age_duration" mapstructure:"max_age_duration"`    // log file max age duration, default is 7 days
	RotationDuration time.Duration `json:"rotation_duration" yaml:"rotation_duration" mapstructure:"rotation_duration"` // log file rotation duration, default is 24h
}

type S3Config struct {
	Region          string `json:"region,omitempty" yaml:"region" mapstructure:"region"`
	AccessKeyID     string `json:"access_key_id,omitempty" yaml:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key,omitempty" yaml:"secret_access_key" mapstructure:"secret_access_key"`
	Bucket          string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	PrefixKey       string `json:"prefix_key" yaml:"prefix_key" mapstructure:"prefix_key"`
}

type RedisConfig struct {
	Network      string `json:"network,omitempty" yaml:"network" mapstructure:"network"`
	Host         string `json:"host,omitempty" yaml:"host" mapstructure:"host"`
	Port         int    `json:"port,omitempty" yaml:"port" mapstructure:"port"`
	Username     string `json:"username,omitempty" yaml:"username" mapstructure:"username"`
	Password     string `json:"password,omitempty" yaml:"password" mapstructure:"password"`
	DB           int    `json:"db,omitempty" yaml:"db" mapstructure:"db"`
	TLSEnabled   bool   `json:"tls_enabled,omitempty" yaml:"tls_enabled" mapstructure:"tls_enabled"`
	Mode         string `json:"mode,omitempty" yaml:"mode" mapstructure:"mode"`
	PoolSize     int    `json:"pool_size,omitempty" yaml:"pool_size" mapstructure:"pool_size"`
	MaxIdleConns int    `json:"max_idle_conns,omitempty" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MinIdleConns int    `json:"min_idle_conns,omitempty" yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
}

type SentryConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	DSN     string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
}

type GrpcServerConfig struct {
	Host string `json:"host" yaml:"host" mapstructure:"host"`
	Port int    `json:"port" yaml:"port" mapstructure:"port"`
}

type ServiceConfig struct {
	Name     string            `json:"name" yaml:"name" mapstructure:"name"`
	Endpoint string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Headers  map[string]string `json:"headers" yaml:"headers" mapstructure:"headers"`
}

func MustGetServiceConfig(svcName ...string) []ServiceConfig {
	res := make([]ServiceConfig, len(svcName))
	for i, svc := range svcName {
		var found bool
		for _, info := range Config.Services {
			if strings.ToLower(info.Name) == strings.ToLower(svc) {
				res[i] = info
				found = true
				break
			}
		}
		if !found {
			panic("config of service " + svc + " not found")
		}
	}
	return res
}

const Environment = "PROFILE"

func ReadConfig(val any) (string, error) {
	filename := "local"
	env, found := os.LookupEnv(Environment)
	if found {
		filename = env
		log.Log.Infof("use %s", env)
	}
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("resources")
	err := v.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Log.Debugw("config file not found", "file", "app.yaml")
		}
	}
	v.AutomaticEnv()
	for _, key := range v.AllKeys() {
		rawValue := v.Get(key)
		val, ok := rawValue.(string)
		if !ok {
			continue
		}
		valueAfter := os.ExpandEnv(val)
		v.Set(key, valueAfter)
	}

	if err := v.Unmarshal(val); err != nil {
		return "", errors.WithMessagef(err, "unmarshal config failed")
	}

	return env, nil
}

func GetOrDefault[T any](v T, def T) T {
	rv := reflect.ValueOf(v)
	switch rv.Type().Kind() {
	case reflect.Struct, reflect.Array, reflect.Chan:
		panic("unsupported type")
	case reflect.Pointer, reflect.Interface:
		if rv.IsNil() {
			return def
		}
		return v
	case reflect.Map, reflect.Slice:
		if rv.IsNil() || rv.Len() == 0 {
			return def
		}
		return v
	default:
		if rv.IsZero() {
			return def
		}
	}
	return v
}
