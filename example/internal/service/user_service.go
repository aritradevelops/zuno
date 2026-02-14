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

type UserFields struct {
	Email string `validate:"required,email"`
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

type UserService interface {
	List(context.Context, *action.Actor, *pagination.Options) (*pagination.Result[*User], error)
	Create(context.Context, *action.Actor, UserFields) (*User, error)
	Update(context.Context, *action.Actor, uuid.UUID, UserFields) (bool, error)
	View(context.Context, *action.Actor, uuid.UUID) (*User, error)
	Delete(context.Context, *action.Actor, uuid.UUID) (bool, error)
	Destroy(context.Context, *action.Actor, uuid.UUID) (bool, error)
	Restore(context.Context, *action.Actor, uuid.UUID) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// List implements [UserService].
func (s *userService) List(ctx context.Context, actor *action.Actor, opts *pagination.Options) (*pagination.Result[*User], error) {
	result, err := s.userRepository.List(ctx, actor, opts)
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(result.Data))
	for idx, user := range result.Data {
		users[idx] = fromRepositoryUser(user)
	}

	converted := &pagination.Result[*User]{
		Data: users,
		Info: result.Info,
	}
	return converted, nil
}

// Create implements [UserService].
func (s *userService) Create(ctx context.Context, actor *action.Actor, payload UserFields) (*User, error) {
	if err := validation.Validate(payload); err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return nil, err
	}
	user, err := s.userRepository.Create(ctx, actor, repository.UserFields(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return nil, err
	}
	logger.Info().Msg("user created successfully!")
	return fromRepositoryUser(user), nil
}

// Update implements [UserService].
func (s *userService) Update(ctx context.Context, actor *action.Actor, id uuid.UUID, payload UserFields) (bool, error) {
	if err := validation.Validate(payload); err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return false, err
	}
	ok, err := s.userRepository.UpdateByID(ctx, actor, id, repository.UserFields(payload))
	if err != nil {
		logger.Error().Err(err).Msg("failed to update user")
		return false, err
	}
	logger.Info().Msg("user updated successfully!")
	return ok, nil
}

// View implements [UserService].
func (s *userService) View(ctx context.Context, actor *action.Actor, id uuid.UUID) (*User, error) {
	user, err := s.userRepository.FindByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("user not found")
		return nil, err
	}
	return fromRepositoryUser(user), nil
}

// Delete implements [UserService].
func (s *userService) Delete(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.userRepository.DeleteByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("user not found")
		return false, err
	}
	logger.Info().Msg("user deleted successfully")
	return ok, nil
}

// Destroy implements [UserService].
func (s *userService) Destroy(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.userRepository.DestroyByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("user not found")
		return false, err
	}
	logger.Info().Msg("user destroyed successfully")
	return ok, nil
}

// Restore implements [UserService].
func (s *userService) Restore(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	ok, err := s.userRepository.RestoreByID(ctx, actor, id)
	if err != nil {
		logger.Error().Err(err).Msg("user not found")
		return false, err
	}
	logger.Info().Msg("user restored successfully")
	return ok, nil
}

func fromRepositoryUser(user *repository.User) *User {
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
