package config

import (
	"fmt"
	"os"
)

const (
	envServerServiceHost     = "SERVER_SERVICE_HOST"
	envServerServicePort     = "SERVER_SERVICE_PORT"
	envServerServiceLogLevel = "SERVER_SERVICE_LOG_LEVEL"

	METRICS_SERVER          = "0.0.0.0:50052"
	METRICS_SERVER_ENDPOINT = "http://0.0.0.0:50052/metrics"
)

// ServerAppConfig ...
type ServerAppConfig struct {
	Host     string
	Port     string
	LogLevel string
}

// GetCombinedAddress with Host and Port
func (cfg *ServerAppConfig) GetCombinedAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

// LoadFromEnv form environment variables
func (cfg *ServerAppConfig) LoadFromEnv() {
	cfg.Host = os.Getenv(envServerServiceHost)
	if len(cfg.Host) == 0 {
		cfg.Host = "0.0.0.0"
	}
	cfg.Port = os.Getenv(envServerServicePort)
	if len(cfg.Port) == 0 {
		cfg.Port = "50051"
	}
	cfg.LogLevel = os.Getenv(envServerServiceLogLevel)
	if len(cfg.LogLevel) == 0 {
		cfg.LogLevel = "local"
	}

}
