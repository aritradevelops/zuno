package add

import (
	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addModulesCmd = &cobra.Command{
	Use:   "modules [name...]",
	Short: "Add modules",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Error("No config found")
			return
		}
		addModules(config, args, cmd)
		logger.Info("Modules added successfully")
	},
}

func addModules(config *config.Config, modules []string, cmd *cobra.Command) {
	addDomains(config, modules, cmd)
	addRepositories(config, modules, cmd)
	addServices(config, modules, cmd)
	addAdapters(config, modules, cmd)
	addTransports(config, modules, cmd)
}

func init() {
	addCmd.AddCommand(addModulesCmd)
}
