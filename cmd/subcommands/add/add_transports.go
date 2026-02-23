package add

import (
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/generators/fiber"

	"github.com/spf13/cobra"
)

var addTransportsCmd = &cobra.Command{
	Use:   "transports [name...]",
	Short: "Add transports",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			cmd.Println("No config found")
			return
		}

		addTransports(config, args, cmd)
	},
}

func addTransports(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		for _, transport := range config.Transports {
			if transport.Type == "http" {
				if transport.Provider == "fiber" {
					if err := fiber.AddNewHandler(config.PackageName, module); err != nil {
						cmd.PrintErrln("failed to add new handler:", err)
						return
					}

					if err := fiber.AddNewRouter(config.PackageName, module); err != nil {
						cmd.PrintErrln("failed to add new router:", err)
						return
					}

					if err := fiber.RegisterNewRouter(module); err != nil {
						cmd.PrintErrln("failed to register new router:", err)
						return
					}
				}
			}
		}
	}
	cmd.Println("Transports added successfully")
}

func init() {
	addCmd.AddCommand(addTransportsCmd)
}
