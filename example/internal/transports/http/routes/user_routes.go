package routes

import (
	"goserve/internal/transports/http/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(router fiber.Router, authMiddleware fiber.Handler, userHandler *handler.UserHandler) {
	router.Get("/list", authMiddleware, userHandler.List)
	router.Post("/create", authMiddleware, userHandler.Create)
	router.Get("/view/:id", authMiddleware, userHandler.View)
	router.Put("/update/:id", authMiddleware, userHandler.Update)
	router.Delete("/delete/:id", authMiddleware, userHandler.Delete)
	router.Delete("/destroy/:id", authMiddleware, userHandler.Destroy)
	router.Delete("/restore/:id", authMiddleware, userHandler.Restore)
}
