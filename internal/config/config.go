package config

import (
	"errors"
	"os"
	"strings"
)

const (
	defaultPort        = "7070"
	defaultServiceName = "notification-service"
)

type Config struct {
	Port        string
	ServiceName string
	DatabaseURL string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return Config{
		Port:        port,
		ServiceName: serviceName,
		DatabaseURL: databaseURL,
	}, nil
}
