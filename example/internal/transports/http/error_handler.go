package http

import (
	"goserve/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, e error) error {
		logger.Error().Err(e).Msg("error occured")
		if err, ok := e.(*fiber.Error); ok {
			c.Status(err.Code)
			return c.JSON(fiber.Map{
				"err":     err.Code,
				"message": err.Message,
			})
		}
		return c.JSON(fiber.Map{
			"err":     e,
			"message": "something went wrong",
		})
	}
}
