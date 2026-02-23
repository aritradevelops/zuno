package service

import (
	"context"
	"goserve/internal/action"
	"goserve/internal/pagination"
	"goserve/internal/repository"
	"goserve/pkg/logger"
	"goserve/pkg/validation"
	"time"

	"github.com/google/uuid"
)

type ProductVariantFields struct {
}

type ProductVariant struct {
	UID uuid.UUID
	ProductVariantFields
	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}

type ProductVariantService interface {
	List(context.Context, *action.Actor, *pagination.Options) (*pagination.Result[*ProductVariant], error)
	Create(context.Context, *action.Actor, ProductVariantFields) (*ProductVariant, error)
	Update(context.Context, *action.Actor, uuid.UUID, ProductVariantFields) (bool, error)
	View(context.Context, *action.Actor, uuid.UUID) (*ProductVariant, error)
	Delete(context.Context, *action.Actor, uuid.UUID) (bool, error)
	Destroy(context.Context, *action.Actor, uuid.UUID) (bool, error)
	Restore(context.Context, *action.Actor, uuid.UUID) (bool, error)
}

type productVariantService struct {
	productVariantRepository repository.ProductVariantRepository
}

func NewProductVariantService(productVariantRepository repository.ProductVariantRepository) ProductVariantService {
	return &productVariantService{
		productVariantRepository: productVariantRepository,
	}
}

// List implements [ProductVariantService].
func (s *productVariantService) List(ctx context.Context, actor *action.Actor, opts *pagination.Options) (*pagination.Result[*ProductVariant], error) {
	result, err := s.productVariantRepository.List(ctx, actor, opts)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list product variants")
		return nil, ConvertRepositoryError(err)
	}
	productVariants := make([]*ProductVariant, len(result.Data))
	for idx, productVariant := range result.Data {
		productVariants[idx] = fromRepositoryProductVariant(productVariant)
	}

	converted := &pagination.Result[*ProductVariant]{
		Data: productVariants,
		Info: result.Info,
	}
	return converted, nil
}

// Create implements [ProductVariantService].
func (s *productVariantService) Create(ctx context.Context, actor *action.Actor, payload ProductVariantFields) (*ProductVariant, error) {
	if err := validation.Validate(payload); err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return nil, NewValidationError("Invalid product variant data provided", err)
	}
	productVariant, err := s.productVariantRepository.Create(ctx, actor, repository.ProductVariantFields(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create product variant")
		return nil, ConvertRepositoryError(err)
	}
	logger.Info().Msg("product variant created successfully!")
	return fromRepositoryProductVariant(productVariant), nil
}

// Update implements [ProductVariantService].
func (s *productVariantService) Update(ctx context.Context, actor *action.Actor, id uuid.UUID, payload ProductVariantFields) (bool, error) {
	if err := validation.Validate(payload); err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return false, NewValidationError("Invalid product variant data provided", err)
	}
	ok, err := s.productVariantRepository.UpdateByID(ctx, actor, id, repository.ProductVariantFields(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to update product variant")
		return false, ConvertRepositoryError(err)
	}
	logger.Info().Msg("product variant updated successfully!")
	return ok, nil
}

// View implements [ProductVariantService].
func (s *productVariantService) View(ctx context.Context, actor *action.Actor, id uuid.UUID) (*ProductVariant, error) {
	productVariant, err := s.productVariantRepository.FindByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("product variant not found")
		return nil, ConvertRepositoryError(err)
	}
	return fromRepositoryProductVariant(productVariant), nil
}

// Delete implements [ProductVariantService].
func (s *productVariantService) Delete(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.productVariantRepository.DeleteByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete product variant")
		return false, ConvertRepositoryError(err)
	}
	logger.Info().Msg("product variant deleted successfully")
	return ok, nil
}

// Destroy implements [ProductVariantService].
func (s *productVariantService) Destroy(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.productVariantRepository.DestroyByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to destroy product variant")
		return false, ConvertRepositoryError(err)
	}
	logger.Info().Msg("product variant destroyed successfully")
	return ok, nil
}

// Restore implements [ProductVariantService].
func (s *productVariantService) Restore(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.productVariantRepository.RestoreByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to restore product variant")
		return false, ConvertRepositoryError(err)
	}
	logger.Info().Msg("product variant restored successfully")
	return ok, nil
}

func fromRepositoryProductVariant(productVariant *repository.ProductVariant) *ProductVariant {
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
