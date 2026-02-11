package add

import (
	"fmt"
	"goserve-cli/pkg/config"
	"goserve-cli/pkg/logger"
	"goserve-cli/pkg/stringx"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

const fieldMarker = "// THIS IS FIELD MARKER. DO NOT TOUCH!"
const mapperMarker = "// THIS IS MAPPER MARKER. DO NOT TOUCH!"

var addFieldCommand = &cobra.Command{
	Use:     "field [module name type]",
	Aliases: []string{"f"},
	Short:   "add a field to repository and its adapter",
	Args:    cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		module := stringx.New(args[0])
		field := stringx.New(args[1])
		fieldType := args[2]
		conf, err := config.Load()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return err
		}
		projectRoot, err := os.Getwd()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return err
		}
		if err = addFieldToRepository(module, field, fieldType, conf, projectRoot); err != nil {
			return err
		}
		if err = addFieldToMongoAdapter(module, field, fieldType, conf, projectRoot); err != nil {
			return err
		}
		return nil
	},
}

func addFieldToRepository(module stringx.Stringx, field stringx.Stringx, fieldType string, conf *config.Config, projectRoot string) error {
	pathToRepositoryFile := path.Join(projectRoot, conf.PathToRepository, module.RepositoryFileName())
	return addFieldToFile(module, field, fieldType, pathToRepositoryFile)
}
func addFieldToMongoAdapter(module stringx.Stringx, field stringx.Stringx, fieldType string, conf *config.Config, projectRoot string) error {
	pathToAdapter := path.Join(projectRoot, "internal", "adapters/mongodb")

	pathToModel := path.Join(pathToAdapter, module.ModelFileName())
	pathToRepositoryAdapter := path.Join(pathToAdapter, module.RepositoryFileName())
	if err := addFieldToFile(module, field, fieldType, pathToModel, "bson"); err != nil {
		return err
	}
	if err := addFieldToFile(module, field, fieldType, pathToRepositoryAdapter); err != nil {
		return err
	}
	return nil
}

func addFieldToFile(module stringx.Stringx, field stringx.Stringx, fieldType string, pathToFile string, tags ...string) error {
	if _, err := os.Stat(pathToFile); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	contentBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return err
	}
	content := string(contentBytes)
	tagStringFinal := ""
	if len(tags) > 0 {
		tagStrings := []string{}
		for _, tag := range tags {
			switch tag {
			case "json":
				tagStrings = append(tagStrings, field.JsonTag())
			case "bson":
				tagStrings = append(tagStrings, field.BsonTag())
			case "validate":
				tagStrings = append(tagStrings, field.ValidateTag())
			}
		}
		tagStringFinal = fmt.Sprintf("`%s`", strings.Join(tagStrings, " "))
	}

	updated := strings.ReplaceAll(content, fieldMarker,
		fmt.Sprintf("%s %s %s \n\t%s", field.ModuleName(), fieldType, tagStringFinal, fieldMarker))

	updated = strings.ReplaceAll(updated, mapperMarker,
		fmt.Sprintf("%s:%s,\n\t%s", field.ModuleName(), fmt.Sprintf("m.%s", field.ModuleName()), mapperMarker))

	if err := os.WriteFile(pathToFile, []byte(updated), 0644); err != nil {
		return err
	}
	return nil
}
