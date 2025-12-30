package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/modelsApp"
)

func CreateEntities(entities []*modelsApp.Entity) {
	convert := make([]interface{}, len(entities))
	for i, entity := range entities {
		convert[i] = entity
	}
	db.InsertMany(CT_MODEL_ENTITIES, convert)
}

func UpdateEntities(entities []*modelsApp.Entity) {
	for _, entity := range entities {
		db.SetEntitiesWhere(CT_MODEL_ENTITIES, bson.M{"_id": entity.Uid}, entity)

	}
}

func CreateConnections(modelUid primitive.ObjectID, connections [][2]primitive.ObjectID) {

	batch := make([]interface{}, len(connections))

	for i, conn := range connections {

		tmp := modelsApp.EntityConnection{}
		tmp.A = conn[0]
		tmp.B = conn[1]
		tmp.IsTwoWay = false
		tmp.Model = modelUid

		batch[i] = tmp
	}

	db.InsertMany(CT_MODEL_CONNECTIONS, batch)
}

func DeleteEntities(entities []primitive.ObjectID) {

	db.DeleteEntities(CT_MODEL_CONNECTIONS, bson.M{"$or": []bson.M{
		{"a": bson.M{"$in": entities}},
		{"b": bson.M{"$in": entities}},
	}})
	db.DeleteEntities(CT_MODEL_ENTITIES, bson.M{"_id": bson.M{"$in": entities}})
}

func DeleteConnection(a, b primitive.ObjectID) {
	db.DeleteEntity(CT_MODEL_CONNECTIONS, bson.M{
		"a": a, "b": b,
	})
}

func ListEntitiesOfWorld(worldUid primitive.ObjectID) []*modelsApp.Entity {
	return QueryWorldEntities(bson.M{"model": worldUid})
}

func ListConnectionsOfWorld(worldUid primitive.ObjectID) []*modelsApp.EntityConnection {
	return db.QueryEntities[modelsApp.EntityConnection](CT_MODEL_CONNECTIONS, bson.M{"model": worldUid})
}

func QueryWorldEntities(query bson.M) []*modelsApp.Entity {
	return db.QueryEntities[modelsApp.Entity](CT_MODEL_ENTITIES, query)
}
