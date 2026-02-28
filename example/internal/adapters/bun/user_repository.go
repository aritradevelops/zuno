package bun

import (
	"context"
	"goserve/internal/action"
	"goserve/internal/pagination"
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Create implements [repository.UserRepository].
func (u *UserRepository) Create(ctx context.Context, actor *action.Actor, payload repository.UserFields) (*repository.User, error) {
	fields := (UserFields)(payload)
	user := &User{
		UID:        uuid.New(),
		UserFields: fields,
		CreatedAt:  time.Now(),
		CreatedBy:  actor.UID,
		UpdatedAt:  time.Now(),
		UpdatedBy:  actor.UID,
	}

	_, err := u.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, repository.NewDatabaseInsertError("create user", err)
	}
	return user.toRepository(), nil
}

// UpdateByID implements [repository.UserRepository].
func (u *UserRepository) UpdateByID(ctx context.Context, actor *action.Actor, id uuid.UUID, payload repository.UserFields) (bool, error) {
	fields := (UserFields)(payload)
	user := &User{
		UID:        id,
		UserFields: fields,
		UpdatedAt:  time.Now(),
		UpdatedBy:  actor.UID,
	}
	filter, args := getFilter(actor)
	_, err := u.db.NewUpdate().
		Model(user).
		OmitZero().
		Where("uid = ?", id).
		Where(filter, args...).
		Exec(ctx)
	if err != nil {
		return false, repository.NewDatabaseUpdateError("update user", err)
	}
	return true, nil
}

// DeleteByID implements [repository.UserRepository].
func (u *UserRepository) DeleteByID(context.Context, *action.Actor, uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// DestroyByID implements [repository.UserRepository].
func (u *UserRepository) DestroyByID(context.Context, *action.Actor, uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// FindByID implements [repository.UserRepository].
func (u *UserRepository) FindByID(context.Context, *action.Actor, uuid.UUID) (*repository.User, error) {
	panic("unimplemented")
}

// List implements [repository.UserRepository].
func (u *UserRepository) List(context.Context, *action.Actor, *pagination.Options) (*pagination.Result[*repository.User], error) {
	panic("unimplemented")
}

// RestoreByID implements [repository.UserRepository].
func (u *UserRepository) RestoreByID(context.Context, *action.Actor, uuid.UUID) (bool, error) {
	panic("unimplemented")
}
