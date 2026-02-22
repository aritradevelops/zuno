
package routes

import (
	"goserve/internal/transports/http/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(router fiber.Router, authMiddleware fiber.Handler, userHandler *handler.UserHandler) {
	router.Get("/", authMiddleware, userHandler.List)
	router.Post("/", authMiddleware, userHandler.Create)
	router.Get("/:id", authMiddleware, userHandler.View)
	router.Put("/:id", authMiddleware, userHandler.Update)
	router.Delete("/:id", authMiddleware, userHandler.Delete)
	router.Delete("/destroy/:id", authMiddleware, userHandler.Destroy)
	router.Delete("/restore/:id", authMiddleware, userHandler.Restore)
}
