package api

import (
	"context"
	"fmt"
	"goserve/internal/adapters/mongodb"
	"goserve/internal/config"
	"goserve/internal/service"
	"goserve/internal/transports/http"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
)

func Run() error {
	// load the configuration
	config, err := config.Load()
	if err != nil {
		return err
	}
	log.Printf("config loaded successfully, %+v", config)
	setupCtx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	db := mongodb.New(config.Database.Connection.Url)
	log.Println("connected to database successfully!")
	if err := db.Connect(setupCtx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	repositories := mongodb.NewRepositories(db.Client(), db.Database())
	services := service.New(repositories)
	s := http.NewServer(config, services)
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	go s.Start()
	signal := <-quitCh
	logger.Info().Str("signal", signal.String()).Msg("shutting down gracefully...")
	if err := s.Shutdown(); err != nil {
		logger.Info().Err(err).Msg("failed to shutdown gracefully!")
		return err
	}
	if err := db.Disconnect(context.Background()); err != nil {
		logger.Info().Err(err).Msg("failed to close database connections gracefully!")
		return err
	}
	return nil
}
