package repository

import (
	"context"
	"goserve/internal/action"
	"goserve/internal/pagination"
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
	List(context.Context, *action.Actor, *pagination.Options) (*pagination.Result[*User], error)
	Create(context.Context, *action.Actor, UserFields) (*User, error)
	FindByID(context.Context, *action.Actor, uuid.UUID) (*User, error)
	UpdateByID(context.Context, *action.Actor, uuid.UUID, UserFields) (bool, error)
	DeleteByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
	DestroyByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
	RestoreByID(context.Context, *action.Actor, uuid.UUID) (bool, error)
}
