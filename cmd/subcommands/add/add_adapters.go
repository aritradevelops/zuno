package add

import (
	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/generators/mongodb"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addAdaptersCmd = &cobra.Command{
	Use:   "adapters [name...]",
	Short: "Add adapters",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Info("No config found")
			return
		}
		addAdapters(config, args, cmd)

	},
}

func addAdapters(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if config.Adapter.Database.Enabled {
			switch config.Adapter.Database.Provider {
			case "mongodb":
				if err := mongodb.AddNewModel(config.Package, module); err != nil {
					logger.Error("failed to add new model:", "err", err)
					return
				}

				if err := mongodb.AddNewRepository(config.Package, module); err != nil {
					logger.Error("failed to add new repository:", "err", err)
					return
				}

				if err := mongodb.RegisterNewRepository(module); err != nil {
					logger.Error("failed to register new repository:", "err", err)
					return
				}
			}
		}
	}
	logger.Info("Adapters added successfully")

}

func init() {
	addCmd.AddCommand(addAdaptersCmd)
}
