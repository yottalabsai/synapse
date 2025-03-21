package config

import (
	"fmt"
	"net/url"
)

const (
	SnowflakeNodeIDRedisKeyForNextId = "synapse:connector:snowflake:nextId"
)

func JoinUrlByServiceName(serviceName string, path string) (string, error) {
	for _, svc := range Config.App.Services {
		if svc.Name == serviceName {
			return url.JoinPath(svc.Endpoint, path)
		}
	}
	return "", fmt.Errorf("client not configured")
}
