package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/generators/mongodb"

	"github.com/spf13/cobra"
)

var addAdaptersCmd = &cobra.Command{
	Use:   "adapters [name...]",
	Short: "Add adapters",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			cmd.Println("No config found")
			return
		}
		addAdapters(config, args, cmd)

		cmd.Println("Adapter added successfully")
	},
}

func addAdapters(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		for _, adapter := range config.Adapters {
			if adapter.Type == "database" {
				if adapter.Provider == "mongodb" {
					if err := mongodb.AddNewModel(config.PackageName, module); err != nil {
						cmd.PrintErrln("failed to add new model:", err)
						return
					}

					if err := mongodb.AddNewRepository(config.PackageName, module); err != nil {
						cmd.PrintErrln("failed to add new repository:", err)
						return
					}

					if err := mongodb.RegisterNewRepository(module); err != nil {
						cmd.PrintErrln("failed to register new repository:", err)
						return
					}
				}
			}
		}
	}
}

func init() {
	addCmd.AddCommand(addAdaptersCmd)
}
