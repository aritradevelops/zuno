package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/generators/service"
	"zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addServicesCmd = &cobra.Command{
	Use:   "services [name...]",
	Short: "Add services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Info("No config found")
			return
		}
		addServices(config, args, cmd)
	},
}

func addServices(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if err := service.AddNewService(config.Package, module); err != nil {
			logger.Error("failed to add new service:", "err", err)
			return
		}

		if err := service.RegisterNewService(module); err != nil {
			logger.Error("failed to register new repository:", "err", err)
			return
		}
	}

	logger.Info("Services added successfully")
}

func init() {
	addCmd.AddCommand(addServicesCmd)
}
