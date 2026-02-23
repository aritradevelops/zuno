package add

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add components to your project",
}

func Cmd() *cobra.Command {
	return addCmd
}
