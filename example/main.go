package main

import (
	"goserve/cmd/api"
	"goserve/pkg/logger"
	"os"
)

func main() {
	if err := api.Run(); err != nil {
		logger.Error().Err(err).Msg("failed to run the api")
		os.Exit(1)
	}
}
