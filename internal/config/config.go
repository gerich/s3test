package config

import (
	"time"

	"github.com/gerich/s3test/internal/adapters/filesystem"
	"github.com/gerich/s3test/internal/ports/rest"
)

type Config struct {
	HTTP             *rest.Config
	FileSystemConfig *filesystem.Config
	LogLevel         string
}

func New() *Config {
	cfg := &Config{HTTP: &rest.Config{}, FileSystemConfig: &filesystem.Config{}}

	// Уйдет в heap
	cfg.HTTP.Port = 8080
	cfg.HTTP.ReadTimeout = 10 * time.Second
	cfg.HTTP.WriteTimeout = 20 * time.Second
	cfg.HTTP.JWTSecret = "9okmnbgvftyuijk"
	cfg.HTTP.AllowedUsers = []string{"foo", "bar", "baz"}
	cfg.HTTP.MaxFileSizeMB = 4

	cfg.FileSystemConfig.Buckets = []string{
		"storage/bucket1",
		"storage/bucket2",
		"storage/bucket3",
		"storage/bucket4",
		"storage/bucket5",
		"storage/bucket6",
		//"storage/bucket7",
	}

	cfg.LogLevel = "debug"

	return cfg
}
