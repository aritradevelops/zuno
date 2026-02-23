package repository

import (
	"context"
	"goserve/internal/action"
	"goserve/internal/pagination"
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

type ProductVariantRepository interface {
	List(context.Context, *action.Actor, *pagination.Options) (*pagination.Result[*ProductVariant], error)
	Create(context.Context, *action.Actor, ProductVariantFields) (*ProductVariant, error)
	FindByID(context.Context, *action.Actor, uuid.UUID) (*ProductVariant, error)
	UpdateByID(context.Context, *action.Actor, uuid.UUID, ProductVariantFields) (bool, error)
	DeleteByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
	DestroyByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
	RestoreByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
}