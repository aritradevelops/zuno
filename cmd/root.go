package cmd

import (
	"github.com/aritradevelops/zuno/cmd/subcommands/add"
	initP "github.com/aritradevelops/zuno/cmd/subcommands/init"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use: "github.com/aritradevelops/zuno",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(add.Cmd())
	rootCmd.AddCommand(initP.Cmd())
}

func Execute() error {
	return rootCmd.Execute()
}
