package handler

import (
	"goserve/internal/pagination"
	"goserve/internal/service"
	"goserve/internal/transports/http/middlewares"
	"goserve/pkg/logger"
	"time"

	"goserve/internal/transports/http/translation"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var UserHandlerInfo = map[string]string{
	"Module": "User",
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type UserFields struct {
	Email string `json:"email" example:"user@example.com"`
}

type User struct {
	UID        uuid.UUID `json:"id" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	UserFields `json:",inline"`
	CreatedAt  time.Time  `json:"created_at" example:"2026-02-12T09:04:30.145+00:00"`
	CreatedBy  uuid.UUID  `json:"created_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	UpdatedAt  time.Time  `json:"updated_at" example:"2026-02-12T09:04:30.145+00:00"`
	UpdatedBy  uuid.UUID  `json:"updated_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	DeletedAt  *time.Time `json:"deleted_at" example:"2026-02-12T09:04:30.145+00:00"`
	DeletedBy  *uuid.UUID `json:"deleted_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
}

// List godoc
// @Summary      list users
// @Description  list users
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        options query PaginationOptions false "options for listing customization"
// @Success      200 {object} Response[[]User, PaginationInfo , NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      403 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/list [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	opts := (*PaginationOptions)(pagination.NewOptions())
	if err := c.Bind().Query(opts); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	result, err := h.userService.List(c.Context(), actor, (*pagination.Options)(opts))
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to list users")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	users := make([]*User, len(result.Data))
	logger.Info().Any("users", users).Msg("users")
	for idx, user := range result.Data {
		users[idx] = fromServiceUser(user)
	}
	response := PaginatedResponse[*User]{
		Data: users,
		Info: PaginationInfo(result.Info),
	}
	return c.JSON(SuccessWithInfo(translation.Localize(c, "controller.list", UserHandlerInfo), response.Data, response.Info))
}

// Create godoc
// @Summary      create an user
// @Description  create an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body UserFields true "create user payload"
// @Success      201 {object} Response[User, NoInfo, NoError]
// @Failure      400 {object} Response[NoData, NoInfo, map[string]validation.ValidationError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      409 {object} Response[NoData, NoInfo, string]
// @Failure      422 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/create [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var payload UserFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	user, err := h.userService.Create(c.Context(), actor, service.UserFields(payload))
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to create user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	c.Status(fiber.StatusCreated)
	return c.JSON(Success(translation.Localize(c, "controller.create", UserHandlerInfo), fromServiceUser(user)))
}

// Update godoc
// @Summary      update an user
// @Description  update an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body UserFields true "update user payload"
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      400 {object} Response[NoData, NoInfo, map[string]validation.ValidationError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      422 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/update/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid user ID"))
	}
	var payload UserFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	ok, err := h.userService.Update(c.Context(), actor, uid, service.UserFields(payload))
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to update user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	return c.JSON(Success(translation.Localize(c, "controller.update", UserHandlerInfo), ok))
}

// View godoc
// @Summary      view an user
// @Description  view an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[User, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/view/{id} [get]
func (h *UserHandler) View(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid user ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	user, err := h.userService.View(c.Context(), actor, uid)
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to view user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	return c.JSON(Success(translation.Localize(c, "controller.view", UserHandlerInfo), fromServiceUser(user)))
}

// Delete godoc
// @Summary      delete an user
// @Description  delete an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/delete/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid user ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	ok, err := h.userService.Delete(c.Context(), actor, uid)
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to delete user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	return c.JSON(Success(translation.Localize(c, "controller.delete", UserHandlerInfo), ok))
}

// Destroy godoc
// @Summary      destroy an user
// @Description  destroy an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/destroy/{id} [delete]
func (h *UserHandler) Destroy(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid user ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	ok, err := h.userService.Destroy(c.Context(), actor, uid)
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to destroy user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	return c.JSON(Success(translation.Localize(c, "controller.destroy", UserHandlerInfo), ok))
}

// Restore godoc
// @Summary      restore an user
// @Description  restore an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /users/restore/{id} [patch]
func (h *UserHandler) Restore(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid user ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return c.Status(fiber.StatusUnauthorized).JSON(Failure(translation.Localize(c, "errors.401"), "Unauthorized"))
	}
	ok, err := h.userService.Restore(c.Context(), actor, uid)
	if err != nil {
		var serviceErr service.Error
		if service.AsServiceError(err, &serviceErr) {
			return err // Return service error as-is for middleware handling
		}
		logger.Error().Err(err).Msg("failed to restore user")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(Failure(translation.Localize(c, "errors.500"), "Internal server error"))
	}
	return c.JSON(Success(translation.Localize(c, "controller.restore", UserHandlerInfo), ok))
}

func fromServiceUser(user *service.User) *User {
	return &User{
		UID: user.UID,
		UserFields: UserFields{
			Email: user.Email,
		},
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
		DeletedAt: user.DeletedAt,
		DeletedBy: user.DeletedBy,
	}
}
