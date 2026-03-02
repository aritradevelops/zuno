package bun

import (
	"path"

	"github.com/aritradevelops/zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewModelData struct {
	Package    string
	FieldsType string
	Module     string
	Table      string
	FileName   string
}

func AddNewModel(packageName, module string) error {
	data, err := prepareAddNewModelData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_model.gotmpl",
		path.Join(pathToModel, data.FileName), data,
	)
}

func prepareAddNewModelData(packageName, module string) (AddNewModelData, error) {
	pluralize := pluralize.NewClient()
	return AddNewModelData{
		Package:    packageName,
		FieldsType: module + "Fields",
		Module:     module,
		Table:      strcase.ToSnake(pluralize.Plural(module)),
		FileName:   strcase.ToSnake(module) + "_model.go",
	}, nil
}
