package knowledge

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"turtle/db"
	"turtle/models"
)

const CT_KNOWLEDGE = "knowledge"

func ListKnowledge(user *models.User) []*Knowledge {

	findOptions := options.FindOptions{}
	findOptions.Projection = bson.M{"typeData": 0}

	return db.QueryEntities[Knowledge](CT_KNOWLEDGE, bson.M{
		"org": user.Org,
	}, &findOptions)
}

func COUKnowledge(user *models.User, knowledge *Knowledge) {
	knowledge.Org = user.Org

	if knowledge.Uid.IsZero() {
		db.InsertEntity(CT_KNOWLEDGE, knowledge)
	} else {
		db.UpdateOneCustom(CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": knowledge})
	}

}

func GetKnowledge(user *models.User, knowledgeUid primitive.ObjectID) *Knowledge {
	return db.QueryEntity[Knowledge](CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}

func DeleteKnowledge(user *models.User, knowledgeUid primitive.ObjectID) {
	db.DeleteEntity(CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}
