package api

import (
	"context"
	"fmt"
	"goserve/internal/action"
	"goserve/internal/adapters/mongodb"
	"goserve/internal/config"
	"goserve/internal/pagination"
	"goserve/internal/repository"
	"log"
	"time"

	"github.com/google/uuid"
)

var actorId = "7e602c5d-8460-4790-b153-a4feb5ceba3a"

func RunTest() error {
	// load the configuration
	config, err := config.Load()
	if err != nil {
		return err
	}
	log.Printf("config loaded successfully, %+v", config)
	setupCtx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	db := mongodb.New(config.Database.Connection.Url)

	if err := db.Connect(setupCtx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Disconnect(context.Background())
	log.Println("connected to database successfully!")
	operationCtx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	actor := &action.Actor{
		UID:   uuid.MustParse(actorId),
		Scope: "owner",
	}
	repo := mongodb.NewUserRepository(db.Client(), db.Database())
	user, err := repo.Create(operationCtx, actor, repository.UserFields{
		Email: "test4@gmail.com",
	})
	if err != nil {
		return fmt.Errorf("failed to created user: %w", err)
	}
	log.Printf("user created successfully, %+v", user)

	ok, err := repo.UpdateByID(operationCtx, actor, user.UID, repository.UserFields{
		Email: "test44@gmail.com",
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	log.Printf("user updated successfully, %+v", ok)

	ok, err = repo.DeleteByID(operationCtx, actor, user.UID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	log.Printf("user deleted successfully, %+v", ok)

	result, err := repo.List(operationCtx, actor, pagination.NewOptions())
	if err != nil {
		return fmt.Errorf("failed to list user: %w", err)
	}

	_, err = repo.RestoreByID(context.TODO(), actor, user.UID)
	if err != nil {
		return fmt.Errorf("failed to list user: %w", err)
	}
	for idx, u := range result.Data {
		fmt.Printf("user %d => %+v\n", idx, u)
	}
	return nil
}
