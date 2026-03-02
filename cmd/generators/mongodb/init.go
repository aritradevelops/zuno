package mongodb

import (
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/utils"
)

func Initialize(config *config.Config) error {
	return utils.CloneTemplates(templates, "templates/init", pathToMongodbAdapter, config)
}
