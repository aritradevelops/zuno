package add

import (
	"goserve-cli/pkg/config"
	"goserve-cli/pkg/logger"
	"goserve-cli/pkg/stringx"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var addRepositoryCmd = &cobra.Command{
	Use:     "repository [module]",
	Aliases: []string{"r", "repo"},
	Short:   "add a repository with its adapter",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		module := stringx.New(args[0])
		conf, err := config.Load()
		if err != nil {
			logger.Error().Err(err).Msg("failed to load config")
			return err
		}
		projectRoot, err := os.Getwd()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return err
		}
		logger.Info().Str("repository", module.ModuleName()).Msg("creating repository")
		if err := createRepositoryInterface(module, conf, projectRoot); err != nil {
			return err
		}
		logger.Info().Msg("creating mongodb adapter")
		if err := createMongoAdapter(module, conf, projectRoot); err != nil {
			return err
		}
		return nil
	},
}

func createRepositoryInterface(module stringx.Stringx, conf *config.Config, projectRoot string) error {

	pathToRepository := path.Join(projectRoot, conf.PathToRepository)
	if _, err := os.Stat(pathToRepository); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(pathToRepository, 0775); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	pathToRepositoryFile := path.Join(pathToRepository, module.RepositoryFileName())
	if _, err := os.Stat(pathToRepositoryFile); err != nil {
		if os.IsExist(err) {
			return err
		}
	}
	f, err := os.Create(pathToRepositoryFile)
	if err != nil {
		return err
	}
	if err := Execute(f, "templates/repository.gotmpl", module); err != nil {
		return err
	}
	return nil
}

func createMongoAdapter(module stringx.Stringx, conf *config.Config, projectRoot string) error {
	pathToAdapter := path.Join(projectRoot, "internal", "adapters/mongodb")

	pathToModel := path.Join(pathToAdapter, module.ModelFileName())
	pathToRepositoryAdapter := path.Join(pathToAdapter, module.RepositoryFileName())

	f, err := os.Create(pathToModel)
	if err != nil {
		return err
	}

	if err := Execute(f, "templates/adapters/mongodb/model.gotmpl", module); err != nil {
		return err
	}

	f, err = os.Create(pathToRepositoryAdapter)
	if err != nil {
		return err
	}

	if err := Execute(f, "templates/adapters/mongodb/repository.gotmpl", module); err != nil {
		return err
	}
	return nil
}
