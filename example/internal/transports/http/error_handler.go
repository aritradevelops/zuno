package http

import (
	"errors"
	"fmt"
	"goserve/internal/service"
	"goserve/internal/transports/http/handler"
	"goserve/internal/transports/http/translation"
	"goserve/pkg/logger"
	"goserve/pkg/validation"
	httpb "net/http"
	"os"
	"strconv"

	"github.com/gobeam/stringy"
	"github.com/gofiber/fiber/v3"
)

func getHTTPStatusForServiceError(err *service.ServiceError) int {
	if err == nil {
		return 200
	}

	switch err.Code {
	case service.ErrAlreadyExists:
		return httpb.StatusConflict
	case service.ErrBadRequest:
		return httpb.StatusBadRequest
	case service.ErrExternalService:
		return httpb.StatusServiceUnavailable
	case service.ErrForbidden:
		return httpb.StatusForbidden
	case service.ErrInternal:
		return httpb.StatusInternalServerError
	case service.ErrNotFound:
		return httpb.StatusNotFound
	case service.ErrOperationBlocked, service.ErrTimeout:
		return httpb.StatusInternalServerError
	case service.ErrUnauthorized:
		return httpb.StatusUnauthorized
	case service.ErrValidationFailed:
		return httpb.StatusUnprocessableEntity
	default:
		return httpb.StatusInternalServerError
	}
}

func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, e error) error {
		logger.Error().Err(e).Msg("error occurred")

		var validationErrors validation.ValidationErrors
		if errors.As(e, &validationErrors) {
			finalError := map[string]validation.ValidationError{}
			for _, err := range validationErrors {
				err.Message = translation.Localize(c, fmt.Sprintf("validation.%s", err.Code), err)
				//! NOTE: mapping error against fields would not work as the field are from service layer
				// and at the moment I can't find a way to map it against the handler layer and get the json
				// tag out of it so our best bet is to just snake casing it and hoping it would work, alternatively
				// we can check against the value for the field which just sounds crazy but it is what it is....
				finalError[stringy.New(err.Field).SnakeCase().ToLower()] = err
			}
			c.Status(httpb.StatusUnprocessableEntity)
			return c.JSON(handler.Failure(translation.Localize(c, "errors.422"), finalError))
		}

		// Handle service errors
		var serviceErr *service.ServiceError
		if service.AsServiceError(e, &serviceErr) {
			httpStatus := getHTTPStatusForServiceError(serviceErr)
			c.Status(httpStatus)
			fmt.Println("serviceErr.Details", serviceErr.Details)
			return c.JSON(handler.Failure(translation.Localize(c, fmt.Sprintf("errors.%d", httpStatus), serviceErr.Details), serviceErr))
		}

		maskedError := fmt.Errorf("something went wrong")
		if isDebug() {
			maskedError = e
		}
		if err, ok := e.(*fiber.Error); ok {
			c.Status(err.Code)
			return c.JSON(handler.Failure(translation.Localize(c, fmt.Sprintf("errors.%d", err.Code)), maskedError))
		}
		c.Status(httpb.StatusInternalServerError)
		return c.JSON(handler.Failure("something went wrong! please try again later", maskedError))
	}
}

func isDebug() bool {
	isItReally, err := strconv.ParseBool((os.Getenv("DEBUG")))
	if err != nil {
		return false
	}
	return isItReally
}
