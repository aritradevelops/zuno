package mongodb

import (
	"html/template"
	"os"
	"path"

	"github.com/ettle/strcase"
)

type AddNewModelData struct {
	Package    string
	FieldsType string
	Module     string
	FileName   string
}

func AddNewModel(packageName, module string) error {
	data, err := prepareAddNewModelData(packageName, module)
	if err != nil {
		return err
	}
	tmplContent, err := loadTemplate("new_model")
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(wd, pathToModel, data.FileName)
	tmpl, err := template.New(filePath).Parse(tmplContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func prepareAddNewModelData(packageName, module string) (AddNewModelData, error) {
	return AddNewModelData{
		Package:    packageName,
		FieldsType: module + "Fields",
		Module:     module,
		FileName:   strcase.ToSnake(module) + "_model.go",
	}, nil
}
