package config

import (
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	RedisKeyPoolRanks = prefixRedisKey("pool:ranks")
)

func InitRedis(cfg *RedisConfig) redis.UniversalClient {
	switch cfg.Mode {
	case "cluster":
		return NewClusterRedisClient(cfg)
	default:
		return NewStandaloneRedisClient(cfg)
	}
}

func NewStandaloneRedisClient(cfg *RedisConfig) *redis.Client {
	options := &redis.Options{
		Network:               cfg.Network,
		Addr:                  fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:              cfg.Username,
		Password:              cfg.Password,
		DB:                    cfg.DB,
		ContextTimeoutEnabled: true,
	}
	if cfg.TLSEnabled {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	return redis.NewClient(options)
}

func NewClusterRedisClient(cfg *RedisConfig) *redis.ClusterClient {
	options := &redis.ClusterOptions{
		Addrs:                 []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Username:              cfg.Username,
		Password:              cfg.Password,
		PoolSize:              20,
		ContextTimeoutEnabled: true,
	}
	if cfg.TLSEnabled {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	return redis.NewClusterClient(options)
}

func prefixRedisKey(key string) string {
	return "exchange:" + key
}
