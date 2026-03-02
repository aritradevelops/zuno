package add

import (
	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/generators/fiber"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addTransportsCmd = &cobra.Command{
	Use:   "transports [name...]",
	Short: "Add transports",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Info("No config found")
			return
		}

		addTransports(config, args, cmd)
	},
}

func addTransports(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if config.Transport.Http.Enabled {
			switch config.Transport.Http.Provider {
			case "fiber":
				if err := fiber.AddNewHandler(config.Package, module); err != nil {
					logger.Error("failed to add new handler:", "err", err)
					return
				}
				if err := fiber.RegisterNewHandler(module); err != nil {
					logger.Error("failed to register new handler:", "err", err)
					return
				}

				if err := fiber.AddNewRouter(config.Package, module); err != nil {
					logger.Error("failed to add new router:", "err", err)
					return
				}

				if err := fiber.RegisterNewRouter(module); err != nil {
					logger.Error("failed to register new router:", "err", err)
					return
				}
			}
		}
	}
	logger.Info("Transports added successfully")
}

func init() {
	addCmd.AddCommand(addTransportsCmd)
}
