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

type ProductVariantRepository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

const ProductVariantCollectionName = "product_variants"

var ProductVariantSearchFields = []string{}

func NewProductVariantRepository(client *mongo.Client, db *mongo.Database) repository.ProductVariantRepository {
	return &ProductVariantRepository{
		client:     client,
		collection: db.Collection(ProductVariantCollectionName),
	}
}

// List implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) List(ctx context.Context, actor *action.Actor, opts *pagination.Options) (*pagination.Result[*repository.ProductVariant], error) {
	productVariants, info, err := paginate[*ProductVariant](ctx, r.collection, actor, opts, ProductVariantSearchFields)
	if err != nil {
		return nil, repository.NewDatabaseQueryError("list product variants", err)
	}
	data := make([]*repository.ProductVariant, len(productVariants))
	for idx, productVariant := range productVariants {
		data[idx] = productVariant.toRepository()
	}
	return &pagination.Result[*repository.ProductVariant]{
		Data: data,
		Info: *info,
	}, nil
}

// Create implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) Create(ctx context.Context, actor *action.Actor, payload repository.ProductVariantFields) (*repository.ProductVariant, error) {
	fields := (ProductVariantFields)(payload)
	productVariant := &ProductVariant{
		ProductVariantFields: fields,
	}
	productVariant.UID = uuid.New()
	productVariant.CreatedAt = time.Now()
	productVariant.CreatedBy = actor.UID
	productVariant.UpdatedAt = time.Now()
	productVariant.UpdatedBy = actor.UID

	if _, err := r.collection.InsertOne(ctx, productVariant); err != nil {
		return nil, repository.NewDatabaseQueryError("create product variant", err)
	}
	return productVariant.toRepository(), nil
}

// FindByID implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) FindByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (*repository.ProductVariant, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, repository.NewNotFoundError("ProductVariant", id.String())
		}
		return nil, repository.NewDatabaseQueryError("find product variant by ID", result.Err())
	}
	var productVariant *ProductVariant
	if err := result.Decode(&productVariant); err != nil {
		return nil, repository.NewDatabaseQueryError("decode product variant", err)
	}
	return productVariant.toRepository(), nil
}

// UpdateByID implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) UpdateByID(ctx context.Context, actor *action.Actor, id uuid.UUID, payload repository.ProductVariantFields) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	dataMap, err := toMap((ProductVariantFields)(payload))
	if err != nil {
		return false, repository.NewInvalidDataError("Invalid product variant data", map[string]any{"field": "payload", "value": payload})
	}
	applyFilter(filter, actor)
	dataMap["updated_at"] = time.Now()
	dataMap["updated_by"] = actor.UID
	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": dataMap})
	if err != nil {
		return false, repository.NewDatabaseQueryError("update product variant", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("ProductVariant", id.String())
	}
	return true, nil
}

// DeleteByID implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) DeleteByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
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
		return false, repository.NewDatabaseQueryError("delete product variant", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("ProductVariant", id.String())
	}
	return true, nil
}

// DestroyByID implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) DestroyByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
	filter := bson.D{
		{Key: "uid", Value: id}, {Key: "deleted_at", Value: nil},
	}
	applyFilter(filter, actor)
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, repository.NewDatabaseQueryError("destroy product variant", err)
	}
	if result.DeletedCount == 0 {
		return false, repository.NewNotFoundError("ProductVariant", id.String())
	}
	return result.Acknowledged, nil
}

// RestoreByID implements repository.ProductVariantRepository.
func (r *ProductVariantRepository) RestoreByID(ctx context.Context, actor *action.Actor, id uuid.UUID) (bool, error) {
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
		return false, repository.NewDatabaseQueryError("restore product variant", err)
	}
	if result.ModifiedCount == 0 {
		return false, repository.NewNotFoundError("ProductVariant", id.String())
	}
	return true, nil
}
