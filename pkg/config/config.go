package config

import (
	"os"
)

type Configuration struct {
	ServiceName             string
	LogLevel                string
	DatabaseUser            string
	DatabaseName            string
	DatabasePassword        string
	ConfigurationServiceURL string
	OtelCollectorEndpoint   string
	OtelCollectorPort       string
}

var Config *Configuration

func LoadConfig() {
	Config = &Configuration{
		ServiceName:             getEnv("SERVICE_NAME", "repositories_service"),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		DatabaseUser:            getEnv("DB_USER", "defaultuser"),
		DatabaseName:            getEnv("DB_NAME", "repositories"),
		DatabasePassword:        getEnv("DB_PASSWORD", "defaultpassword"),
		ConfigurationServiceURL: getEnv("CONFIGURATION_SERVICE_URL", "http://localhost:8081"),
		OtelCollectorEndpoint:   getEnv("OTEL_COLLECTOR_ENDPOINT", "localhost"),
		OtelCollectorPort:       getEnv("OTEL_COLLECTOR_PORT", "4317"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
