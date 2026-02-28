package bun

import (
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserFields struct {
	Email string  `bun:"email,unique,notnull"`
	Name  string  `bun:"name,notnull,type:varchar(255)"`
	Dp    *string `bun:"dp,type:varchar(255),nullzero"`
}

type User struct {
	bun.BaseModel `bun:"table:users"`
	UID           uuid.UUID `bun:"uid,pk,type:uuid"`
	UserFields    `bun:",inline"`
	CreatedAt     time.Time  `bun:"created_at,nullzero,notnull,type:timestamptz"`
	CreatedBy     uuid.UUID  `bun:"created_by,nullzero,notnull,type:uuid"`
	UpdatedAt     time.Time  `bun:"updated_at,nullzero,notnull,type:timestamptz"`
	UpdatedBy     uuid.UUID  `bun:"updated_by,nullzero,notnull,type:uuid"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero,type:timestamptz"`
	DeletedBy     *uuid.UUID `bun:"deleted_by,nullzero,type:uuid"`
}

func (u *User) toRepository() *repository.User {
	return &repository.User{
		UID: u.UID,
		UserFields: repository.UserFields{
			Email: u.Email, Name: u.Name,
			Dp: u.Dp,
		},
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
		DeletedAt: u.DeletedAt,
		DeletedBy: u.DeletedBy,
	}
}
