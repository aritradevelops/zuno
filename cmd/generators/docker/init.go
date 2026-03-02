package docker

import (
	"path"
	"zuno/cmd/config"
	"zuno/cmd/utils"
)

func Initialize(config *config.Config, verbose bool) error {
	switch config.Adapter.Database.Provider {
	case "bun":
		return initPostgresDocker(config, verbose)
	}
	return nil
}

func initPostgresDocker(config *config.Config, verbose bool) error {
	return utils.CloneTemplates(templates, "templates/postgres", path.Join(pathToDocker, "postgres"), config)
}
