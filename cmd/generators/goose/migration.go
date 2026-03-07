package goose

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aritradevelops/zuno/cmd/data"
	"github.com/aritradevelops/zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewMigrationData struct {
	Package string
	Table   string
}

type AddNewColumnsMigrationData struct {
	Table  string
	Fields []data.Field
}

func AddNewCreateTableMigration(packageName, module string, pathToMigration string) error {
	data, err := prepareAddNewMigrationData(packageName, module)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_create_table_migration.gotmpl",
		path.Join(pathToMigration, getCreateTableMigrationName(pathToMigration, data.Table)), data,
	)
}

func AddNewColumnsMigration(packageName, module string, pathToMigration string, fields []data.Field) error {
	data, err := prepareAddNewColumnsMigrationData(packageName, module, fields)
	if err != nil {
		return err
	}
	return utils.CreateFromTemplate(
		templates, "templates/new_add_columns_migration.gotmpl",
		path.Join(pathToMigration, getAddNewColumnsMigrationName(pathToMigration, data.Table, data.Fields)), data,
	)
}

func getCreateTableMigrationName(pathToMigration, tableName string) string {
	max, _ := getMaxMigrationNumber(pathToMigration)
	return fmt.Sprintf("%05d_create_%s_table.sql", max+1, tableName)
}
func getAddNewColumnsMigrationName(pathToMigration, tableName string, fields []data.Field) string {
	max, _ := getMaxMigrationNumber(pathToMigration)
	columnNames := []string{}
	for _, field := range fields {
		columnNames = append(columnNames, strcase.ToSnake(field.Name))
	}
	return fmt.Sprintf("%05d_add_%s_to_%s.sql", max+1, strings.Join(columnNames, "_"), tableName)
}

func getMaxMigrationNumber(pathToMigrations string) (int, error) {
	files, err := os.ReadDir(pathToMigrations)
	if err != nil {
		return 0, err
	}

	max := 0

	for _, file := range files {
		first := strings.SplitN(file.Name(), "_", 2)[0]
		firstInt, _ := strconv.Atoi(first)
		if firstInt > max {
			max = firstInt
		}
	}

	return max, nil
}

func prepareAddNewMigrationData(packageName, module string) (AddNewMigrationData, error) {
	pluralize := pluralize.NewClient()
	return AddNewMigrationData{
		Package: packageName,
		Table:   strcase.ToSnake(pluralize.Plural(module)),
	}, nil
}

func prepareAddNewColumnsMigrationData(packageName, module string, fields []data.Field) (AddNewColumnsMigrationData, error) {
	pluralize := pluralize.NewClient()

	return AddNewColumnsMigrationData{
		Table:  strcase.ToSnake(pluralize.Plural(module)),
		Fields: fields,
	}, nil
}
