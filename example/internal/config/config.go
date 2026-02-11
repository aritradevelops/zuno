package config

import (
	// automatically loads the environment variables from .env file

	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

var k = koanf.New("_")

type Config struct {
	Database Database
	Timeout  time.Duration
}

type Database struct {
	Connection Connection
}

type Connection struct {
	Url string `validate:"required,url"`
}

func Load() (*Config, error) {
	if err := k.Load(env.Provider("", "_", func(s string) string {
		return strings.ToTitle(strings.ToLower(s))
	}), nil); err != nil {
		return nil, err
	}
	defaults := &Config{
		Timeout: time.Minute * 2,
	}
	if err := k.Unmarshal("", defaults); err != nil {
		return nil, err
	}
	// TODO: validation
	return defaults, nil
}
