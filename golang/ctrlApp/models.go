package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"turtle/db"
	"turtle/modelsApp"
)

const CT_MODELS = "models"
const CT_MODEL_ENTITIES = "model_entities"
const CT_MODEL_CONNECTIONS = "model_connections"

func COUModel(ct *modelsApp.Model) {
	if ct.Uid.IsZero() {
		db.InsertEntity(CT_MODELS, ct)
	} else {
		db.UpdateEntity(CT_MODELS, ct)
	}
}

func DeleteModel(uid primitive.ObjectID) {

	DeleteActorsOfModel(uid)

	db.DeleteEntities(CT_MODEL_ENTITIES, bson.M{
		"model": uid,
	})

	db.DeleteEntities(CT_MODEL_CONNECTIONS, bson.M{
		"model": uid,
	})

	db.DeleteEntity(CT_MODELS, bson.M{
		"_id": uid,
	})
}

func ListModels() []modelsApp.Model {
	return db.QueryEntitiesAsCopy[modelsApp.Model](CT_MODELS, bson.M{})
}

func GetModel(uid primitive.ObjectID) *modelsApp.Model {
	return db.QueryEntity[modelsApp.Model](CT_MODELS, bson.M{"_id": uid})
}
