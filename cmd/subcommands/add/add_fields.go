package add

import (
	"fmt"
	"strings"
	"zuno/cmd/app"
	"zuno/cmd/config"
	"zuno/cmd/data"
	"zuno/cmd/generators/fiber"
	"zuno/cmd/generators/mongodb"
	"zuno/cmd/generators/repository"
	"zuno/cmd/generators/service"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var addFieldsCmd = &cobra.Command{
	Use:   "fields [module]",
	Short: "Add fields",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Ctx.Config
		if config == nil {
			cmd.Println("No config found")
			return
		}

		var fields []data.Field
		addAnother := true

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
						Options(
							huh.NewOption("string", "string"),
							huh.NewOption("int", "int"),
							huh.NewOption("float64", "float64"),
							huh.NewOption("bool", "bool"),
							huh.NewOption("bson.ObjectID", "bson.ObjectID"),
							huh.NewOption("*bson.ObjectID", "*bson.ObjectID"),
						).
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
		cmd.Println("Fields added successfully")
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
		cmd.Println("Error adding fields:", err)
		return
	}

	err = service.AddFieldsToService(module, fields)
	if err != nil {
		cmd.Println("Error adding fields:", err)
		return
	}

	err = mongodb.AddFieldsToModel(module, fields)
	if err != nil {
		cmd.Println("Error adding fields:", err)
		return
	}

	err = fiber.AddFieldsToHandler(module, fields)
	if err != nil {
		cmd.Println("Error adding fields:", err)
		return
	}

}

func init() {
	addCmd.AddCommand(addFieldsCmd)
}
