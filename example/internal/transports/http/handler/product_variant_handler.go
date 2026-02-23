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

var ProductVariantHandlerInfo = map[string]string{
	"Module": "ProductVariant",
}

type ProductVariantHandler struct {
	productVariantService service.ProductVariantService
}

func NewProductVariantHandler(productVariantService service.ProductVariantService) *ProductVariantHandler {
	return &ProductVariantHandler{
		productVariantService: productVariantService,
	}
}

type ProductVariantFields struct {
	
}

type ProductVariant struct {
	UID        uuid.UUID `json:"id" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	ProductVariantFields `json:",inline"`
	CreatedAt  time.Time  `json:"created_at" example:"2026-02-12T09:04:30.145+00:00"`
	CreatedBy  uuid.UUID  `json:"created_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	UpdatedAt  time.Time  `json:"updated_at" example:"2026-02-12T09:04:30.145+00:00"`
	UpdatedBy  uuid.UUID  `json:"updated_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
	DeletedAt  *time.Time `json:"deleted_at" example:"2026-02-12T09:04:30.145+00:00"`
	DeletedBy  *uuid.UUID `json:"deleted_by" example:"7e602c5d-8460-4790-b153-a4feb5ceba3a"`
}

// List godoc
// @Summary      list product variants
// @Description  list product variants
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param        options query PaginationOptions false "options for listing customization"
// @Success      200 {object} Response[[]ProductVariant, PaginationInfo , NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      403 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/list [get]
func (h *ProductVariantHandler) List(c fiber.Ctx) error {
	opts := (*PaginationOptions)(pagination.NewOptions())
	if err := c.Bind().Query(opts); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		logger.Error().Err(err).Msg("can't get actor")
		return fiber.ErrUnauthorized
	}
	result, err := h.productVariantService.List(c.Context(), actor, (*pagination.Options)(opts))
	if err != nil {
		return err
	}
	productVariants := make([]*ProductVariant, len(result.Data))
	for idx, productVariant := range result.Data {
		productVariants[idx] = fromServiceProductVariant(productVariant)
	}
	response := PaginatedResponse[*ProductVariant]{
		Data: productVariants,
		Info: PaginationInfo(result.Info),
	}
	return c.JSON(SuccessWithInfo(translation.Localize(c, "controller.list", ProductVariantHandlerInfo), response.Data, response.Info))
}

// Create godoc
// @Summary      create a product variant
// @Description  create a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param        request body ProductVariantFields true "create product variant payload"
// @Success      201 {object} Response[ProductVariant, NoInfo, NoError]
// @Failure      400 {object} Response[NoData, NoInfo, map[string]validation.ValidationError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      409 {object} Response[NoData, NoInfo, string]
// @Failure      422 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/create [post]
func (h *ProductVariantHandler) Create(c fiber.Ctx) error {
	var payload ProductVariantFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	productVariant, err := h.productVariantService.Create(c.Context(), actor, service.ProductVariantFields(payload))
	if err != nil {
		return err
	}
	c.Status(fiber.StatusCreated)
	return c.JSON(Success(translation.Localize(c, "controller.create", ProductVariantHandlerInfo), fromServiceProductVariant(productVariant)))
}

// Update godoc
// @Summary      update a product variant
// @Description  update a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param        request body ProductVariantFields true "update product variant payload"
// @Param 			 id path string true "id of the product variant"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      400 {object} Response[NoData, NoInfo, map[string]validation.ValidationError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      422 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/update/{id} [put]
func (h *ProductVariantHandler) Update(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid product variant ID"))
	}
	var payload ProductVariantFields
	if err := c.Bind().Body(&payload); err != nil {
		return err
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	ok, err := h.productVariantService.Update(c.Context(), actor, uid, service.ProductVariantFields(payload))
	if err != nil {
		return err
	}
	return c.JSON(Success(translation.Localize(c, "controller.update", ProductVariantHandlerInfo), ok))
}

// View godoc
// @Summary      view a product variant
// @Description  view a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the product variant"
// @Success      200 {object} Response[ProductVariant, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/view/{id} [get]
func (h *ProductVariantHandler) View(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid product variant ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	productVariant, err := h.productVariantService.View(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success(translation.Localize(c, "controller.view", ProductVariantHandlerInfo), fromServiceProductVariant(productVariant)))
}

// Delete godoc
// @Summary      delete a product variant
// @Description  delete a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the product variant"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/delete/{id} [delete]
func (h *ProductVariantHandler) Delete(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid product variant ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	ok, err := h.productVariantService.Delete(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success(translation.Localize(c, "controller.delete", ProductVariantHandlerInfo), ok))
}

// Destroy godoc
// @Summary      destroy a product variant
// @Description  destroy a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the product variant"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/destroy/{id} [delete]
func (h *ProductVariantHandler) Destroy(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid product variant ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	ok, err := h.productVariantService.Destroy(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success(translation.Localize(c, "controller.destroy", ProductVariantHandlerInfo), ok))
}

// Restore godoc
// @Summary      restore a product variant
// @Description  restore a product variant
// @Tags        product-variant
// @Accept       json
// @Produce      json
// @Param 			 id path string true "id of the product variant"
// @Success      200 {object} Response[bool, NoInfo, NoError]
// @Failure      401 {object} Response[NoData, NoInfo, string]
// @Failure      404 {object} Response[NoData, NoInfo, string]
// @Failure      500 {object} Response[NoData, NoInfo, string]
// @Router       /product-variants/restore/{id} [patch]
func (h *ProductVariantHandler) Restore(c fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Failure(translation.Localize(c, "errors.400"), "Invalid product variant ID"))
	}
	actor, err := middlewares.GetActor(c)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	ok, err := h.productVariantService.Restore(c.Context(), actor, uid)
	if err != nil {
		return err
	}
	return c.JSON(Success(translation.Localize(c, "controller.restore", ProductVariantHandlerInfo), ok))
}

func fromServiceProductVariant(productVariant *service.ProductVariant) *ProductVariant {
	return &ProductVariant{
		UID: productVariant.UID,
		ProductVariantFields: ProductVariantFields{
		},
		CreatedAt: productVariant.CreatedAt,
		CreatedBy: productVariant.CreatedBy,
		UpdatedAt: productVariant.UpdatedAt,
		UpdatedBy: productVariant.UpdatedBy,
		DeletedAt: productVariant.DeletedAt,
		DeletedBy: productVariant.DeletedBy,
	}
}
