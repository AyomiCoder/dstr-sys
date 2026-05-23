package config

import "os"

const (
	defaultPort        = "7070"
	defaultServiceName = "notification-service"
)

type Config struct {
	Port        string
	ServiceName string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	return Config{
		Port:        port,
		ServiceName: serviceName,
	}
}
