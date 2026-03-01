package goose

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"zuno/cmd/utils"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
)

type AddNewMigrationData struct {
	Package string
	Table   string
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

func getCreateTableMigrationName(pathToMigration, tableName string) string {
	max, _ := getMaxMigrationNumber(pathToMigration)
	return fmt.Sprintf("%05d_create_%s_table.sql", max+1, tableName)
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
