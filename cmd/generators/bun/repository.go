package bun

import (
	"path"
	"strings"

	"github.com/aritradevelops/zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewRepositoryData struct {
	Package        string
	RepositoryType string
	Module         string
	Variable       string
	VariablePlural string
	FieldsType     string
	Readable       string
	ReadablePlural string
	FileName       string
}

func AddNewRepository(packageName, module string) error {
	data, err := prepareAddNewRepositoryData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_repository.gotmpl",
		path.Join(pathToRepository, data.FileName), data,
	)
}

func prepareAddNewRepositoryData(packageName, module string) (AddNewRepositoryData, error) {
	pluralize := pluralize.NewClient()
	return AddNewRepositoryData{
		Package:        packageName,
		RepositoryType: module + "Repository",
		Module:         module,
		Variable:       strcase.ToCamel(module),
		VariablePlural: pluralize.Plural(strcase.ToCamel(module)),
		FieldsType:     module + "Fields",
		Readable:       strings.ReplaceAll(strcase.ToKebab(module), "-", " "),
		ReadablePlural: pluralize.Plural(strings.ReplaceAll(strcase.ToKebab(module), "-", " ")),
		FileName:       strcase.ToSnake(module) + "_repository.go",
	}, nil
}
