package mongodb

import (
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ProductVariantFields struct {

}

type ProductVariant struct {
	UID        uuid.UUID `bson:"uid,omitempty"`
	ProductVariantFields `bson:",inline"`
	CreatedAt  time.Time  `bson:"created_at,omitempty"`
	CreatedBy  uuid.UUID  `bson:"created_by,omitempty"`
	UpdatedAt  time.Time  `bson:"updated_at,omitempty"`
	UpdatedBy  uuid.UUID  `bson:"updated_by,omitempty"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
	DeletedBy  *uuid.UUID `bson:"deleted_by,omitempty"`
}

func (u *ProductVariant) toRepository() *repository.ProductVariant {
	return &repository.ProductVariant{
		UID: u.UID,
		ProductVariantFields: repository.ProductVariantFields{
		},
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
		DeletedAt: u.DeletedAt,
		DeletedBy: u.DeletedBy,
	}
}
