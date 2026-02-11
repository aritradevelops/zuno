package mongodb

import (
	"goserve/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func applyFilter(filter bson.D, actor *repository.Actor) {
	switch actor.Scope {
	case "owner":
		filter = append(filter, bson.E{Key: "created_by", Value: actor.UID})
	}
}
