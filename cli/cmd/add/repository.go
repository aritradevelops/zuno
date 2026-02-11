package add

import (
	"fmt"
	"goserve-cli/pkg/config"
	"goserve-cli/pkg/logger"
	"os"
	"path"
	"strings"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

var pluralizer = pluralize.NewClient()

var addRepositoryCmd = &cobra.Command{
	Use:     "repository [name]",
	Aliases: []string{"r", "repo"},
	Short:   "add a repository with its adapter",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := pluralizer.Singular((strings.ToLower(args[0])))
		vars := map[string]string{
			"ModuleTitle":       strings.ToTitle(name),
			"ModuleLowerPlural": pluralizer.Plural(name),
			"ModuleLower":       name,
		}
		logger.Info().Str("repository", name).Msg("creating repository")
		return createRepositoryInterface(name, vars)
	},
}

func createRepositoryInterface(name string, vars map[string]string) error {
	conf, err := config.Load()
	tmpl, err := template.New("root").ParseFS(templateFs, "templates/repository.gotmpl")
	if err != nil {
		return err
	}
	projectRoot, err := os.Getwd()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get working directory")
		return err
	}
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
	pathToRepositoryFile := path.Join(pathToRepository, fmt.Sprintf("%s_repository.go", name))
	if _, err := os.Stat(pathToRepositoryFile); err != nil {
		if os.IsExist(err) {
			return err
		}
	}
	f, err := os.Create(pathToRepositoryFile)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(f, "repository.gotmpl", vars)
}
