package config

type Config struct {
	PathToRepository string
	Adapters         []string
}

func Load() (*Config, error) {
	return &Config{
		PathToRepository: "internal/repository",
		Adapters:         []string{"mongodb"},
	}, nil
}
