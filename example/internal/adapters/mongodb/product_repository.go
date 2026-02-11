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

type ProductRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

const ProductCollectionName = "products"

var ProductSearchFields = []string{}

func NewProductRepository(client *mongo.Client, db *mongo.Database) repository.ProductRepository {
	return &ProductRepository{
		client:     client,
		collection: db.Collection(ProductCollectionName),
	}
}

// List implements repository.ProductRepository.
func (r *ProductRepository) List(ctx context.Context, actor *repository.Actor, opts *repository.ListOptions) (*repository.ListResponse[*repository.Product], error) {
	filter := bson.D{
		{Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	if opts.Search != "" && len(ProductSearchFields) > 0 {
		for _, field := range ProductSearchFields {
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
	var products []*Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	var data []*repository.Product
	for _, product := range products {
		data = append(data, product.toRepository())
	}
	return &repository.ListResponse[*repository.Product]{
		Data: data,
	}, nil
}

// Create implements repository.ProductRepository.
func (r *ProductRepository) Create(ctx context.Context, actor *repository.Actor, payload repository.ProductFields) (*repository.Product, error) {
	fields := (ProductFields)(payload)
	product := &Product{
		ProductFields: fields,
	}
	product.UID = uuid.New()
	product.CreatedAt = time.Now()
	product.CreatedBy = actor.UID
	product.UpdatedAt = time.Now()
	product.UpdatedBy = actor.UID

	if _, err := r.collection.InsertOne(ctx, product); err != nil {
		return nil, err
	}
	return product.toRepository(), nil
}

// FindByID implements repository.ProductRepository.
func (r *ProductRepository) FindByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (*repository.Product, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var product *Product
	if err := result.Decode(&product); err != nil {
		return nil, err
	}
	return product.toRepository(), nil
}

// UpdateByID implements repository.ProductRepository.
func (r *ProductRepository) UpdateByID(ctx context.Context, actor *repository.Actor, id uuid.UUID, payload repository.ProductFields) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	dataMap, err := (ProductFields)(payload).toMap()
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

// DeleteByID implements repository.ProductRepository.
func (r *ProductRepository) DeleteByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (bool, error) {
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

func (r *ProductRepository) DestroyByID(ctx context.Context, actor *repository.Actor, id uuid.UUID) (bool, error) {
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
