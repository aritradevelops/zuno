package fiber

import (
	"zuno/cmd/config"
	"zuno/cmd/utils"
)

func Initialize(cfg *config.Config, verbose bool) error {
	if err := utils.CloneTemplates(templates, "templates/init", pathToHttpProvider, cfg); err != nil {
		return err
	}
	if err := utils.RunCmd("go", verbose, "install", "github.com/swaggo/swag/cmd/swag@latest"); err != nil {
		return err
	}

	if err := utils.RunCmd("swag", verbose, "init", "-g", "internal/transports/http/routes/register.go"); err != nil {
		return err
	}
	return nil
}
