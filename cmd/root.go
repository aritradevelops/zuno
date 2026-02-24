package cmd

import (
	"zuno/cmd/subcommands/add"
	initP "zuno/cmd/subcommands/init"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use: "zuno",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(add.Cmd())
	rootCmd.AddCommand(initP.Cmd())
}

func Execute() error {
	return rootCmd.Execute()
}
