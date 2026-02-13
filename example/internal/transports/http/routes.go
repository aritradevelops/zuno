package http

import (
	"goserve/internal/transports/http/middlewares"
)

func (s *Server) setupRoutes() {
	authMiddleware := middlewares.AuthMiddleware()
	s.app.Post("/api/v1/users/create", authMiddleware, s.handlers.User.Create)
}
