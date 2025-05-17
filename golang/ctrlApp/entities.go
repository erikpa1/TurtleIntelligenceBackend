package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/modelsApp"
)

const CT_ENTITIES = "entities"

func CreateEntities(entities []*modelsApp.Entity) {
	convert := make([]interface{}, len(entities))
	for i, entity := range entities {
		convert[i] = entity
	}
	db.InsertMany(CT_ENTITIES, convert)
}

func UpdateEntities(entities []*modelsApp.Entity) {
	for _, entity := range entities {
		db.UpdateEntitiesWhere(CT_ENTITIES, bson.M{"_id": entity.Uid}, entity)

	}
}

func DeleteEntities(entities []primitive.ObjectID) {
	db.DeleteEntities(CT_ENTITIES, bson.M{"_id": bson.M{"$in": entities}})
}

func ListEntitiesOfWorld(worldUid primitive.ObjectID) []*modelsApp.Entity {
	return QueryWorldEntities(bson.M{"model": worldUid})
}

func QueryWorldEntities(query bson.M) []*modelsApp.Entity {
	return db.QueryEntities[modelsApp.Entity](CT_ENTITIES, query)
}
