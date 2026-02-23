package routes

import (
	"goserve/internal/transports/http/handler"
	"goserve/internal/transports/http/middlewares"

	"github.com/gofiber/fiber/v3"
)

func Register(app fiber.Router, handlers *handler.Handlers) {
	authMiddleware := middlewares.AuthMiddleware()
	RegisterUserRoutes(app.Group("/users"), authMiddleware, handlers.User)
	RegisterProductVariantRoutes(app.Group("/product-variants"), authMiddleware, handlers.ProductVariant)
}
