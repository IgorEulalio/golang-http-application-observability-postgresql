package config

import (
	"fmt"
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
	RabbitMqUsername        string
	RabbitMqPassword        string
	RabbitMqHost            string
	RabbitMqPort            string
}

var Config *Configuration

func LoadConfig() []error {

	var errs []error

	Config = &Configuration{
		ServiceName:             getEnv("SERVICE_NAME", ""),
		LogLevel:                getEnv("LOG_LEVEL", "info"),
		DatabaseUser:            getEnv("DB_USER", "defaultuser"),
		DatabaseName:            getEnv("DB_NAME", "repositories"),
		DatabasePassword:        getEnv("DB_PASSWORD", "defaultpassword"),
		ConfigurationServiceURL: getEnv("CONFIGURATION_SERVICE_URL", "http://localhost:8081"),
		OtelCollectorEndpoint:   getEnv("OTEL_COLLECTOR_ENDPOINT", "localhost"),
		OtelCollectorPort:       getEnv("OTEL_COLLECTOR_PORT", "4317"),
		RabbitMqUsername:        getEnv("RABBITMQ_USERNAME", "guest"),
		RabbitMqPassword:        getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMqHost:            getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMqPort:            getEnv("RABBITMQ_PORT", "5672"),
	}

	if Config.ServiceName == "" {
		errs = append(errs, fmt.Errorf("SERVICE_NAME is missing"))
	}

	return errs
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
