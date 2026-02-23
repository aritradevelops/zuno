package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/generators/service"

	"github.com/spf13/cobra"
)

var addServicesCmd = &cobra.Command{
	Use:   "services [name...]",
	Short: "Add services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			cmd.Println("No config found")
			return
		}
		addServices(config, args, cmd)
	},
}

func addServices(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if err := service.AddNewService(config.PackageName, module); err != nil {
			cmd.PrintErrln("failed to add new service:", err)
			return
		}

		if err := service.RegisterNewService(module); err != nil {
			cmd.PrintErrln("failed to register new repository:", err)
			return
		}
	}

	cmd.Println("Services added successfully")
}

func init() {
	addCmd.AddCommand(addServicesCmd)
}
