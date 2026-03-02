package add

import (
	"os"

	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var configPath string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add components to your project",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			logger.Error("failed to load config:", "err", err)
			os.Exit(1)
		}
		app.Ctx.Config = cfg
	},
}

func init() {
	addCmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"zuno.yml",
		"Path to config file",
	)
}

func Cmd() *cobra.Command {
	return addCmd
}
