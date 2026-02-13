package middlewares

import (
	"fmt"
	"goserve/internal/action"
	"goserve/internal/auth"
	"goserve/pkg/logger"
	"strings"

	"github.com/gofiber/fiber/v3"
)

const authCtxKey = "actor"

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// try getting the token from cookies
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			// try getting from Authorization header
			bearerToken := c.Get("Authorization")
			if bearerToken == "" {
				// not authorized
				return fiber.ErrUnauthorized
			}
			accessToken = strings.TrimLeft(bearerToken, "Bearer ")
		}
		actor, err := auth.Authenticate(accessToken)
		if err != nil {
			logger.Error().Err(err).Msg("invalid access token")
			return fiber.ErrUnauthorized
		}

		parts := strings.Split(c.Route().Path, "/")

		if err = auth.Authorize(parts[3], parts[4], actor); err != nil {
			logger.Error().Err(err).Msg("not having permission for this route")
			return fiber.ErrForbidden
		}
		c.Locals(authCtxKey, actor)
		return c.Next()
	}
}

func GetActor(c fiber.Ctx) (*action.Actor, error) {
	payload := c.Locals(authCtxKey)
	if payload == nil {
		return nil, fmt.Errorf("Actor is only available for protected routes")
	}

	actor, ok := payload.(*action.Actor)
	if !ok {
		return nil, fmt.Errorf("Actor is only available for protected routes")
	}

	return actor, nil
}
