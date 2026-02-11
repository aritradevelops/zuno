package mongodb

import (
	"goserve/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserFields struct {
	Email string `bson:"email,omitempty"`
}

type User struct {
	UID        uuid.UUID `bson:"uid,omitempty"`
	UserFields `bson:",inline"`
	CreatedAt  time.Time  `bson:"created_at,omitempty"`
	CreatedBy  uuid.UUID  `bson:"created_by,omitempty"`
	UpdatedAt  time.Time  `bson:"updated_at,omitempty"`
	UpdatedBy  uuid.UUID  `bson:"updated_by,omitempty"`
	DeletedAt  *time.Time `bson:"deleted_at,omitempty"`
	DeletedBy  *uuid.UUID `bson:"deleted_by,omitempty"`
}

func (u *User) toRepository() *repository.User {
	return &repository.User{
		UID: u.UID,
		UserFields: repository.UserFields{
			Email: u.Email,
		},
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
		DeletedAt: u.DeletedAt,
		DeletedBy: u.DeletedBy,
	}
}

func (u UserFields) toMap() (map[string]any, error) {
	raw, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	if err := bson.Unmarshal(raw, data); err != nil {
		return nil, err
	}
	return data, nil
}
