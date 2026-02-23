package mongodb

import (
	"context"
	"goserve/internal/action"
	"goserve/internal/pagination"
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

const UserCollectionName = "users"

var UserSearchFields = []string{"email"}

func NewUserRepository(client *mongo.Client, db *mongo.Database) repository.UserRepository {
	return &UserRepository{
		client:     client,
		collection: db.Collection(UserCollectionName),
	}
}

// List implements repository.UserRepository.
func (r *UserRepository) List(ctx context.Context, actor *action.Actor, opts *pagination.Options) (*pagination.Result[*repository.User], error) {
	users, info, err := paginate[*User](ctx, r.collection, actor, opts)
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

// Create implements repository.UserRepository.
func (r *UserRepository) Create(ctx context.Context, actor *action.Actor, payload repository.UserFields) (*repository.User, error) {
	fields := (UserFields)(payload)
	user := &User{
		UserFields: fields,
	}
	user.UID = uuid.New()
	user.CreatedAt = time.Now()
	user.CreatedBy = actor.UID
	user.UpdatedAt = time.Now()
	user.UpdatedBy = actor.UID

	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		return nil, repository.NewDatabaseInsertError("create user", err)
	}
	return user.toRepository(), nil
}

// FindByID implements repository.UserRepository.
func (r *UserRepository) FindByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (*repository.User, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, repository.NewNotFoundError("User", id.String())
		}
		return nil, repository.NewDatabaseQueryError("find user by ID", result.Err())
	}
	var user *User
	if err := result.Decode(&user); err != nil {
		return nil, repository.NewDatabaseQueryError("decode user", err)
	}
	return user.toRepository(), nil
}

// UpdateByID implements repository.UserRepository.
func (r *UserRepository) UpdateByID(ctx context.Context, actor *action.Actor, id uuid.UUID, payload repository.UserFields) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	dataMap, err := toMap((UserFields)(payload))
	if err != nil {
		return false, repository.NewInvalidDataError("Invalid user data", map[string]any{"field": "payload", "value": payload})
	}
	applyFilter(filter, actor)
	dataMap["updated_at"] = time.Now()
	dataMap["updated_by"] = actor.UID
	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": dataMap})
	if err != nil {
		return false, repository.NewDatabaseUpdateError("update user", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("User", id.String())
	}
	return true, nil
}

// DeleteByID implements repository.UserRepository.
func (r *UserRepository) DeleteByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"updated_at": time.Now(),
		"updated_by": actor.UID,
		"deleted_at": time.Now(),
		"deleted_by": actor.UID,
	}})
	if err != nil {
		return false, repository.NewDatabaseUpdateError("delete user", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("User", id.String())
	}
	return true, nil
}

// DestroyByID implements repository.UserRepository.
func (r *UserRepository) DestroyByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, repository.NewDatabaseDeleteError("destroy user", err)
	}
	if result.DeletedCount == 0 {
		return false, repository.NewNotFoundError("User", id.String())
	}
	return result.Acknowledged, nil
}

// RestoreByID implements repository.UserRepository.
func (r *UserRepository) RestoreByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: bson.M{"$ne": nil}},
	}
	applyFilter(filter, actor)
	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"updated_at": time.Now(),
		"updated_by": actor.UID,
		"deleted_at": nil,
		"deleted_by": actor.UID,
	}})
	if err != nil {
		return false, repository.NewDatabaseUpdateError("restore user", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("User", id.String())
	}
	return true, nil
}
