package add

import (
	"goserve-cli/pkg/logger"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info().Msg("add called")
	},
}

func init() {
	addCmd.AddCommand(addRepositoryCmd)
}

func Cmd() *cobra.Command {
	return addCmd
}
