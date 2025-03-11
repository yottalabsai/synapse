package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func InitDatasource(_ context.Context, zapLogger *zap.Logger, datasourceCfg *DatasourceConfig) error {
	gormCfg := datasourceCfg.GORM

	cfg := &gorm.Config{
		PrepareStmt: true,
		Logger:      buildLogger(gormCfg.LogLevel, zapLogger),
	}

	var err error
	db, err := gorm.Open(postgres.Open(buildDSN(datasourceCfg.Postgres)), cfg)
	if err != nil {
		return errors.WithMessagef(err, "failed to open database connection")
	}

	sDB, err := db.DB()
	if err != nil {
		return errors.WithMessagef(err, "failed to get sql db")
	}
	sDB.SetMaxOpenConns(datasourceCfg.Pool.MaxOpenConns)
	sDB.SetMaxIdleConns(datasourceCfg.Pool.MaxIdleConns)
	sDB.SetConnMaxLifetime(datasourceCfg.Pool.ConnMaxLifetime)
	sDB.SetConnMaxIdleTime(datasourceCfg.Pool.ConnMaxIdleTime)

	DB = db
	
	return nil
}

func buildLogger(levelStr string, zapLogger *zap.Logger) logger.Interface {
	levelStr = strings.ToLower(levelStr)
	level := logger.Silent
	switch levelStr {
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	}
	l := zapgorm2.New(zapLogger)
	l.IgnoreRecordNotFoundError = true
	return l.LogMode(level)
}

func buildDSN(cfg PostgresConfig) string {
	var dsn strings.Builder
	dsn.WriteString("host=")
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	dsn.WriteString(cfg.Host)

	dsn.WriteString(" port=")
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	dsn.WriteString(fmt.Sprintf("%d", cfg.Port))

	dsn.WriteString(" user=")
	dsn.WriteString(cfg.Username)

	if cfg.Password != "" {
		dsn.WriteString(" password=")
		dsn.WriteString(cfg.Password)
	}

	if cfg.Database != "" {
		dsn.WriteString(" dbname=")
		dsn.WriteString(cfg.Database)
	}

	if cfg.SearchPath != "" {
		dsn.WriteString(" search_path=")
		dsn.WriteString(cfg.SearchPath)
	}

	if cfg.SslMode != "" {
		dsn.WriteString(" sslmode=")
		dsn.WriteString(cfg.SslMode)
	}

	if cfg.TimeZone != "" {
		dsn.WriteString(" TimeZone=")
		dsn.WriteString(cfg.TimeZone)
	}
	return dsn.String()
}
