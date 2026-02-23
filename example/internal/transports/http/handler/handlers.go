package handler

import "goserve/internal/service"

type Handlers struct {
	User           *UserHandler
	ProductVariant *ProductVariantHandler
}

func New(services *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(services.User), ProductVariant: NewProductVariantHandler(services.ProductVariant),
	}
}
