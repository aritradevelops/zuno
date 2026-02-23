package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	PackageName string      `koanf:"package_name"`
	Adapters    []Adapter   `koanf:"adapters"`
	Transports  []Transport `koanf:"transports"`
}

type Transport struct {
	Type     string `koanf:"type"`
	Provider string `koanf:"provider"`
}

type Adapter struct {
	Type     string `koanf:"type"`
	Provider string `koanf:"provider"`
}

func Load(path string) (*Config, error) {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	// Optional env override
	_ = k.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
