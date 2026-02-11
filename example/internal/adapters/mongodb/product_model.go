package mongodb

import (
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductFields struct {
	Name string `bson:"name"`
	// THIS IS FIELD MARKER. DO NOT TOUCH!
}

type Product struct {
	UID           uuid.UUID `bson:"uid,omitempty"`
	ProductFields `bson:",inline"`
	CreatedAt     time.Time  `bson:"created_at,omitempty"`
	CreatedBy     uuid.UUID  `bson:"created_by,omitempty"`
	UpdatedAt     time.Time  `bson:"updated_at,omitempty"`
	UpdatedBy     uuid.UUID  `bson:"updated_by,omitempty"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty"`
	DeletedBy     *uuid.UUID `bson:"deleted_by,omitempty"`
}

func (m *Product) toRepository() *repository.Product {
	return &repository.Product{
		UID: m.UID,
		ProductFields: repository.ProductFields{
			Name: m.Name,
			// THIS IS MAPPER MARKER. DO NOT TOUCH!
		},
		CreatedAt: m.CreatedAt,
		CreatedBy: m.CreatedBy,
		UpdatedAt: m.UpdatedAt,
		UpdatedBy: m.UpdatedBy,
		DeletedAt: m.DeletedAt,
		DeletedBy: m.DeletedBy,
	}
}

func (m ProductFields) toMap() (map[string]any, error) {
	raw, err := bson.Marshal(m)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	if err := bson.Unmarshal(raw, data); err != nil {
		return nil, err
	}
	return data, nil
}
