package handler

import "goserve/internal/service"

type Handlers struct {
	User *UserHandler
}

func New(services *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(services.User),
	}
}
