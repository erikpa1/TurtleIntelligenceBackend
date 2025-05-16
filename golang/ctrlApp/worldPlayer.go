package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
)

func PlayWorld(worldUid primitive.ObjectID) {

	entities := QueryWorldEntities(bson.M{"model": worldUid})

	for _, entity := range entities {
		lg.LogI(entity)
	}

}
