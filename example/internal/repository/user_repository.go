package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserFields struct {
	Email string
}

type User struct {
	UID uuid.UUID
	UserFields
	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
}

type UserRepository interface {
	List(context.Context, *Actor, *ListOptions) (*ListResponse[*User], error)
	Create(context.Context, *Actor, UserFields) (*User, error)
	FindByID(context.Context, *Actor, uuid.UUID) (*User, error)
	UpdateByID(context.Context, *Actor, uuid.UUID, UserFields) (bool, error)
	DeleteByID(context.Context, *Actor, uuid.UUID) (bool, error)
	DestroyByID(context.Context, *Actor, uuid.UUID) (bool, error)
}
