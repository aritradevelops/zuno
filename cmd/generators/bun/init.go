package bun

import (
	"zuno/cmd/config"
	"zuno/cmd/utils"
)

func Initialize(config *config.Config) error {
	if err := utils.CloneTemplates(templates, "templates/init", pathToBunAdapter, config); err != nil {
		return err
	}
	return nil
}
