package mongodb

import (
	"zuno/cmd/config"
	"zuno/cmd/utils"
)

func Initialize(config *config.Config) error {
	return utils.CloneTemplates(templates, "templates/init", pathToMongodbAdapter, config)
}
