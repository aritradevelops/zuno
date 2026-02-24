package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/generators/repository"
	"zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addRepositoriesCmd = &cobra.Command{
	Use:   "repositories [name...]",
	Short: "Add repositories",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Error("No config found")
			return
		}
		addRepositories(config, args, cmd)
	},
}

func addRepositories(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if err := repository.AddNewRepository(config.Package, module); err != nil {
			logger.Error("failed to add new repository:", "err", err)
			return
		}

		if err := repository.RegisterNewRepository(module); err != nil {
			logger.Error("failed to register new repository:", "err", err)
			return
		}
	}

	logger.Info("Repositories added successfully")
}

func init() {
	addCmd.AddCommand(addRepositoriesCmd)
}
