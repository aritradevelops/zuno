package bun

import (
	"context"
	"database/sql"
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

// List implements [repository.UserRepository].
func (u *UserRepository) List(ctx context.Context, actor *action.Actor, opts *pagination.Options) (*pagination.Result[*repository.User], error) {
	users, info, err := paginate[*User](ctx, u.db, actor, opts, UserSearchFields)
	if err != nil {
		return nil, repository.NewDatabaseQueryError("list users", err)
	}
	data := make([]*repository.User, len(users))
	for idx, user := range users {
		data[idx] = user.toRepository()
	}
	return &pagination.Result[*repository.User]{
		Data: data,
		Info: *info,
	}, nil
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
		return nil, repository.NewDatabaseQueryError("create user", err)
	}
	return user.toRepository(), nil
}

// FindByID implements [repository.UserRepository].
func (u *UserRepository) FindByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (*repository.User, error) {
	user := &User{}
	err := u.db.NewSelect().
		Model(user).
		Where("uid = ?", id).
		Where("deleted_at IS NULL").
		ApplyQueryBuilder(scopedFilter(actor)).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NewNotFoundError("user", id.String())
		}
		return nil, repository.NewDatabaseQueryError("find user", err)
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
	cols := append(fields.Fields(), "updated_at", "updated_by")
	result, err := u.db.NewUpdate().
		Model(user).
		Column(cols...).
		Where("uid = ?", id).
		Where("deleted_at IS NULL").
		ApplyQueryBuilder(scopedFilter(actor)).
		Exec(ctx)
	if err != nil {
		return false, repository.NewDatabaseQueryError("update user", err)
	}

	if count, err := result.RowsAffected(); err != nil {
		return false, repository.NewDatabaseQueryError("update user", err)
	} else if count == 0 {
		return false, repository.NewNotFoundError("user", id.String())
	}
	return true, nil
}

// DeleteByID implements [repository.UserRepository].
func (u *UserRepository) DeleteByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	result, err := u.db.NewUpdate().
		Model(&User{}).
		Set("deleted_at = ?", time.Now()).
		Set("deleted_by = ?", actor.UID).
		Where("uid = ?", id).
		Where("deleted_at IS NULL").
		ApplyQueryBuilder(scopedFilter(actor)).
		Exec(ctx)
	if err != nil {
		return false, repository.NewDatabaseQueryError("delete user", err)
	}

	if count, err := result.RowsAffected(); err != nil {
		return false, repository.NewDatabaseQueryError("delete user", err)
	} else if count == 0 {
		return false, repository.NewNotFoundError("user", id.String())
	}
	return true, nil
}

// DestroyByID implements [repository.UserRepository].
func (u *UserRepository) DestroyByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	result, err := u.db.NewDelete().
		Model(&User{}).
		Where("uid = ?", id).
		ApplyQueryBuilder(scopedFilter(actor)).
		Exec(ctx)
	if err != nil {
		return false, repository.NewDatabaseQueryError("delete user", err)
	}

	if count, err := result.RowsAffected(); err != nil {
		return false, repository.NewDatabaseQueryError("delete user", err)
	} else if count == 0 {
		return false, repository.NewNotFoundError("user", id.String())
	}
	return true, nil
}

// RestoreByID implements [repository.UserRepository].
func (u *UserRepository) RestoreByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	result, err := u.db.NewUpdate().
		Model(&User{}).
		Set("deleted_at = NULL").
		Where("uid = ?", id).
		Where("deleted_at IS NOT NULL").
		ApplyQueryBuilder(scopedFilter(actor)).
		Exec(ctx)
	if err != nil {
		return false, repository.NewDatabaseQueryError("restore user", err)
	}

	if count, err := result.RowsAffected(); err != nil {
		return false, repository.NewDatabaseQueryError("restore user", err)
	} else if count == 0 {
		return false, repository.NewNotFoundError("user", id.String())
	}
	return true, nil
}
