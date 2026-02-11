package mongodb

import (
	"context"
	"goserve/internal/repository"
	"regexp"
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

var UserSearchFields = []string{}

func NewUserRepository(client *mongo.Client, db *mongo.Database) repository.UserRepository {
	return &UserRepository{
		client:     client,
		collection: db.Collection(UserCollectionName),
	}
}

// List implements repository.UserRepository.
func (r *UserRepository) List(ctx context.Context, actor *repository.Actor, opts *repository.ListOptions) (*repository.ListResponse[*repository.User], error) {
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	if opts.Search != "" && len(UserSearchFields) > 0 {
		for _, field := range UserSearchFields {
			filter = append(filter, bson.E{Key: field, Value: bson.D{
				{Key: "$regex", Value: regexp.MustCompile(opts.Search)},
			}})
		}
	}
	matchStage := bson.D{
		{Key: "$match", Value: filter},
	}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage})
	if err != nil {
		return nil, err
	}
	var users []*User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	var data []*repository.User
	for _, u := range users {
		data = append(data, u.toRepository())
	}
	return &repository.ListResponse[*repository.User]{
		Data: data,
	}, nil
}

// Create implements repository.UserRepository.
func (r *UserRepository) Create(ctx context.Context, actor *repository.Actor, payload repository.UserFields) (*repository.User, error) {
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
		return nil, err
	}
	return user.toRepository(), nil
}

// FindByID implements repository.UserRepository.
func (r *UserRepository) FindByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (*repository.User, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var user *User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return user.toRepository(), nil
}

// UpdateByID implements repository.UserRepository.
func (r *UserRepository) UpdateByID(ctx context.Context, actor *repository.Actor, id uuid.UUID, payload repository.UserFields) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	dataMap, err := (UserFields)(payload).toMap()
	if err != nil {
		return false, err
	}
	applyFilter(filter, actor)
	dataMap["updated_at"] = time.Now()
	dataMap["updated_by"] = actor.UID
	if _, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": dataMap}); err != nil {
		return false, err
	}
	return true, nil
}

// DeleteByID implements repository.UserRepository.
func (r *UserRepository) DeleteByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	if _, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"updated_at": time.Now(),
		"updated_by": actor.UID,
		"deleted_at": time.Now(),
		"deleted_by": actor.UID,
	}}); err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) DestroyByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	return result.Acknowledged, nil
}
