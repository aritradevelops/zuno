package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductFields struct {
	Name string  
	// THIS IS FIELD MARKER. DO NOT TOUCH!
}

type Product struct {
	UID uuid.UUID
	ProductFields
	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}

type ProductRepository interface {
	List(context.Context, *Actor, *ListOptions) (*ListResponse[*Product], error)
	Create(context.Context, *Actor, ProductFields) (*Product, error)
	FindByID(context.Context, *Actor, uuid.UUID) (*Product, error)
	UpdateByID(context.Context, *Actor, uuid.UUID, ProductFields) (bool, error)
	DeleteByID(context.Context, *Actor, uuid.UUID) (bool, error)
	DestroyByID(context.Context, *Actor, uuid.UUID) (bool, error)
}

