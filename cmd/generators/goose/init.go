package goose

import (
	"os"
	"zuno/cmd/config"
	"zuno/cmd/utils"
)

func Initialize(config *config.Config, verbose bool) error {

	if err := installGoose(verbose); err != nil {
		return err
	}

	if err := addGooseConfigToEnv(config); err != nil {
		return err
	}

	return nil
}

func installGoose(verbose bool) error {
	return utils.RunCmd("go", verbose, "install", "github.com/pressly/goose/v3/cmd/goose@latest")
}

func addGooseConfigToEnv(config *config.Config) error {
	content, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	content = append(content, []byte(`
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:admin@localhost:5432/`+config.PackageBase+`
GOOSE_MIGRATION_DIR=./internal/adapters/bun/migrations
	`)...)

	return os.WriteFile(".env", content, 0644)
}
