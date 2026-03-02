package add

import (
	"fmt"
	"strings"

	"github.com/aritradevelops/zuno/cmd/app"
	"github.com/aritradevelops/zuno/cmd/config"
	"github.com/aritradevelops/zuno/cmd/data"
	"github.com/aritradevelops/zuno/cmd/generators/fiber"
	"github.com/aritradevelops/zuno/cmd/generators/mongodb"
	"github.com/aritradevelops/zuno/cmd/generators/repository"
	"github.com/aritradevelops/zuno/cmd/generators/service"
	"github.com/aritradevelops/zuno/pkg/logger"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var fieldTypes = []string{
	"string",
	"int",
	"float64",
	"bool",
	"bson.ObjectID",
	"*bson.ObjectID",
}

var addFieldsCmd = &cobra.Command{
	Use:   "fields [module]",
	Short: "Add fields",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			logger.Info("No config found")
			return
		}

		var fields []data.Field
		addAnother := true
		options := make([]huh.Option[string], 0)
		for _, fieldType := range fieldTypes {
			options = append(options, huh.NewOption(fieldType, fieldType))
		}
		for addAnother {
			var field data.Field

			preview := formatFieldsPreview(fields)

			if err := huh.NewForm(
				huh.NewGroup(
					huh.NewNote().
						Title("Current Session Fields").
						Description(preview),

					huh.NewInput().
						Title("data.Field name").
						Value(&field.Name),
					huh.NewSelect[string]().
						Title("data.Field type").
						Options(options...).
						Value(&field.Type),
				),
			).Run(); err != nil {
				return
			}

			fields = append(fields, field)

			if err := huh.NewConfirm().
				Title("Add another field?").
				Value(&addAnother).
				Run(); err != nil {
				addAnother = false
			}
		}

		addFields(config, args[0], fields, cmd)
		logger.Info("Fields added successfully")
	},
}

func formatFieldsPreview(fields []data.Field) string {
	if len(fields) == 0 {
		return "No fields added yet."
	}

	var b strings.Builder
	b.WriteString("Fields to be added:\n\n")

	for i, f := range fields {
		b.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, f.Name, f.Type))
	}

	return b.String()
}
func addFields(config *config.Config, module string, fields []data.Field, cmd *cobra.Command) {
	err := repository.AddFieldsToRepository(module, fields)
	if err != nil {
		logger.Error("Error adding fields:", "err", err)
		return
	}

	err = service.AddFieldsToService(module, fields)
	if err != nil {
		logger.Error("Error adding fields:", "err", err)
		return
	}

	err = mongodb.AddFieldsToModel(module, fields)
	if err != nil {
		logger.Error("Error adding fields:", "err", err)
		return
	}

	err = fiber.AddFieldsToHandler(module, fields)
	if err != nil {
		logger.Error("Error adding fields:", "err", err)
		return
	}

}

func init() {
	addCmd.AddCommand(addFieldsCmd)
}
