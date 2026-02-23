package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"

	"github.com/spf13/cobra"
)

var addModuleCmd = &cobra.Command{
	Use:   "modules [name...]",
	Short: "Add modules",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			cmd.Println("No config found")
			return
		}
		addModules(config, args, cmd)
		cmd.Println("Modules added successfully")
	},
}

func addModules(config *config.Config, modules []string, cmd *cobra.Command) {
	addRepositories(config, modules, cmd)
	addServices(config, modules, cmd)
	addAdapters(config, modules, cmd)
	addTransports(config, modules, cmd)
}

func init() {
	addCmd.AddCommand(addModuleCmd)
}
