package cmd

import (
	"log"
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/subcommands/add"

	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use: "zuno",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		app.Ctx.Config = cfg
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"zuno.yml",
		"Path to config file",
	)

	rootCmd.AddCommand(add.Cmd())
}

func Execute() error {
	return rootCmd.Execute()
}
