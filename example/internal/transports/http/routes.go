package http

import (
	"goserve/internal/transports/http/middlewares"
)

func (s *Server) setupRoutes() {
	authMiddleware := middlewares.AuthMiddleware()
	s.app.Get("/api/v1/users/list", authMiddleware, s.handlers.User.List)
	s.app.Post("/api/v1/users/create", authMiddleware, s.handlers.User.Create)
	s.app.Put("/api/v1/users/update/:id", authMiddleware, s.handlers.User.Update)
	s.app.Get("/api/v1/users/view/:id", authMiddleware, s.handlers.User.View)
	s.app.Delete("/api/v1/users/delete/:id", authMiddleware, s.handlers.User.Delete)
	s.app.Delete("/api/v1/users/destroy/:id", authMiddleware, s.handlers.User.Destroy)
	s.app.Patch("/api/v1/users/restore/:id", authMiddleware, s.handlers.User.Restore)
}
