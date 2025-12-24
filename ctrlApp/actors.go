package ctrlApp

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/modelsApp"
)

const CT_ACTORS = "actors"

func QueryActors(query bson.M) []*modelsApp.Actor {
	return db.QueryEntities[modelsApp.Actor](CT_ACTORS, query)
}

func GetActor(uid primitive.ObjectID) *modelsApp.Actor {
	return db.QueryEntity[modelsApp.Actor](CT_ACTORS, bson.M{"_id": uid})
}

func COUActor(actor *modelsApp.Actor) {
	if actor.Uid.IsZero() {
		db.InsertEntity(CT_ACTORS, actor)
	} else {
		db.UpdateOneCustom(CT_ACTORS, bson.M{"_id": actor.Uid}, bson.M{"$set": actor})
	}

}

func DeleteActor(uid primitive.ObjectID) {
	db.DeleteEntity(CT_ACTORS, bson.M{"_id": uid})
}

func DeleteActorsOfModel(uid primitive.ObjectID) {
	db.DeleteEntity(CT_ACTORS, bson.M{"model": uid})
}
