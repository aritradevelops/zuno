package bun

import (
	"goserve/internal/repository"

	"github.com/uptrace/bun"
)

func NewRepositories(db *bun.DB) *repository.Repositories {
	return &repository.Repositories{
		User: NewUserRepository(db),
	}
}
