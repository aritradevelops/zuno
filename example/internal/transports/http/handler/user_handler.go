package handler

import (
	"goserve/internal/pagination"
	"goserve/internal/service"
	"goserve/internal/transports/http/middlewares"
	"goserve/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

var e = fiber.Error{}

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
// @Router       /users/list [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	opts := (*PaginationOptions)(pagination.NewOptions())
	if err := c.Bind().Query(opts); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	result, err := h.userService.List(c.Context(), actor, (*pagination.Options)(opts))
	if err != nil {
		return err
	}
	var users []*User
	for _, user := range result.Data {
		users = append(users, fromServiceUser(user))
	}
	response := PaginatedResponse[*User]{
		Data: users,
		Info: PaginationInfo(result.Info),
	}
	return c.JSON(SuccessWithInfo("controller.list", response.Data, response.Info))
}

// Create godoc
// @Summary      create an user
// @Description  create an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body UserFields true "create user payload"
// @Success      200 {object} Response[User, NoInfo, NoError]
// @Router       /users/create [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var payload UserFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	user, err := h.userService.Create(c.Context(), actor, service.UserFields(payload))
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.create", fromServiceUser(user)))
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
// @Router       /users/update/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	var payload UserFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	ok, err := h.userService.Update(c.Context(), actor, uid, service.UserFields(payload))
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.update", ok))
}

// View godoc
// @Summary      view an user
// @Description  view an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[User, NoInfo, NoError]
// @Router       /users/view/{id} [get]
func (h *UserHandler) View(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	user, err := h.userService.View(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.view", fromServiceUser(user)))
}

// Delete godoc
// @Summary      delete an user
// @Description  delete an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Router       /users/delete/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	ok, err := h.userService.Delete(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.delete", ok))
}

// Destroy godoc
// @Summary      destroy an user
// @Description  destroy an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Router       /users/destroy/{id} [delete]
func (h *UserHandler) Destroy(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	ok, err := h.userService.Destroy(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.destroy", ok))
}

// Restore godoc
// @Summary      restore an user
// @Description  restore an user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the user"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Router       /users/restore/{id} [patch]
func (h *UserHandler) Restore(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	ok, err := h.userService.Restore(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success("controller.restore", ok))
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
