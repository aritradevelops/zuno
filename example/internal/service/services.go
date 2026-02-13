package service

import "goserve/internal/repository"

type Services struct {
	User UserService
}

func New(repositories *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(repositories.User),
	}
}
