package ctrlApp

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_ACTORS = "actors"

func QueryActors(query bson.M) []*modelsApp.Actor {
	return db.QueryEntities[modelsApp.Actor](CT_ACTORS, query)
}

func GetActor(uid primitive.ObjectID) *modelsApp.Actor {
	return db.QueryEntity[modelsApp.Actor](CT_ACTORS, bson.M{"_id": uid})
}

func CreateActor(ct *modelsApp.Actor) {
	db.InsertEntity(CT_ACTORS, ct)
}

func UpdateActor(ct *modelsApp.Actor) {
	db.UpdateOneCustom(CT_ACTORS, bson.M{"_id": ct.Uid}, bson.M{"$set": ct})
}

func DeleteActor(uid primitive.ObjectID) {
	db.DeleteEntity(CT_ACTORS, bson.M{"_id": uid})
}

func DeleteActorsOfModel(uid primitive.ObjectID) {
	db.DeleteEntity(CT_ACTORS, bson.M{"model": uid})
}
