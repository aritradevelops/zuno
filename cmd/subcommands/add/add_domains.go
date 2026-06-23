package add

import (
	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/generators/domain"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/spf13/cobra"
)

var addDomainsCmd = &cobra.Command{
	Use:   "domains [name...]",
	Short: "Add domains",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Info("No config found")
			return
		}
		addServices(config, args, cmd)
	},
}

func addDomains(config *config.Config, modules []string, cmd *cobra.Command) {
	for _, module := range modules {
		if err := domain.AddNewDomain(config.Package, module); err != nil {
			logger.Error("failed to add new domain:", "err", err)
			return
		}
	}

	logger.Info("Domains added successfully")
}

func init() {
	addCmd.AddCommand(addDomainsCmd)
}
