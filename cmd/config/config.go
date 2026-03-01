package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	Package     string    `koanf:"package_name"`
	PackageBase string    `koanf:"package_base"`
	Adapter     Adapter   `koanf:"adapter"`
	Transport   Transport `koanf:"transport"`
}

type Transport struct {
	Http HttpTransport `koanf:"http"`
	Grpc GrpcTransport `koanf:"grpc"`
	Ws   WsTransport   `koanf:"ws"`
}

type HttpTransport struct {
	Enabled  bool   `koanf:"enabled"`
	Provider string `koanf:"provider"`
}

type GrpcTransport struct {
	Enabled  bool   `koanf:"enabled"`
	Provider string `koanf:"provider"`
}

type WsTransport struct {
	Enabled  bool   `koanf:"enabled"`
	Provider string `koanf:"provider"`
}

type Adapter struct {
	Database DatabaseAdapter `koanf:"database"`
}

type DatabaseAdapter struct {
	Enabled   bool            `koanf:"enabled"`
	Provider  string          `koanf:"provider"`
	Migration MigrationConfig `koanf:"migration"`
}

type MigrationConfig struct {
	Enabled  bool   `koanf:"enabled"`
	Provider string `koanf:"provider"`
}

func (c *Config) ToYaml() ([]byte, error) {
	k.Load(structs.Provider(c, "koanf"), nil)
	return k.Marshal(yaml.Parser())
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
