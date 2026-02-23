package routes

import (
	"goserve/internal/transports/http/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterProductVariantRoutes(router fiber.Router, authMiddleware fiber.Handler, productVariantHandler *handler.ProductVariantHandler) {
	router.Get("/", authMiddleware, productVariantHandler.List)
	router.Post("/", authMiddleware, productVariantHandler.Create)
	router.Get("/:id", authMiddleware, productVariantHandler.View)
	router.Put("/:id", authMiddleware, productVariantHandler.Update)
	router.Delete("/:id", authMiddleware, productVariantHandler.Delete)
	router.Delete("/destroy/:id", authMiddleware, productVariantHandler.Destroy)
	router.Delete("/restore/:id", authMiddleware, productVariantHandler.Restore)
}