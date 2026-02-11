package mongodb

import (
	"context"
	"fmt"
	"goserve/internal/adapters/mongodb"
	"goserve/internal/config"
	"goserve/pkg/logger"
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
)

var db *mongodb.MongoDB

func setup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	envPath := path.Join(wd, "../../../../.env.test")
	// load the environment variables
	if err := godotenv.Load(envPath); err != nil {
		logger.Error().Err(err).Msgf("failed to load env file: %s for tests", envPath)
		return err
	}

	// load configuration
	conf, err := config.Load()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load config")
		return err
	}

	// connect to the database
	db = mongodb.New(conf.Database.Connection.Url)
	if err := db.Connect(context.Background()); err != nil {
		logger.Error().Err(err).Msg("failed to connect to the database")
	}
	return nil
}

func teardown() error {
	if db == nil {
		err := fmt.Errorf("database is not initialized")
		logger.Error().Err(err).Msg("failed to disconnect the database")
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	// setup
	if err := setup(); err != nil {
		os.Exit(1)
	}
	defer teardown()
	exitCode := m.Run()
	os.Exit(exitCode)
}
