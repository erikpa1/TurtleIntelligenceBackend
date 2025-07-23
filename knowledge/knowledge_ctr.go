package knowledge

import (
	"context"
	"github.com/erikpa1/turtle/db"
	"github.com/erikpa1/turtle/llm/llmCtrl"
	"github.com/erikpa1/turtle/llm/llmModels"
	"github.com/erikpa1/turtle/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_KNOWLEDGE = "knowledge"
const CT_KNOWLEDGE_EMBEDDING = "knowledge_embedding"

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
		knowledge.Uid = primitive.NewObjectID()

		db.InsertEntity(CT_KNOWLEDGE, knowledge)
	} else {
		db.UpdateOneCustom(CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": knowledge})
	}

	knowledgeTypeData := KnowledgePlainTextTypeData(knowledge.TypeData)
	embeddableText := knowledge.Name + " " + knowledge.Description + " " + knowledgeTypeData.GetEmbeddableString()
	embedding, embError := llmCtrl.CreateStringEmbedding(context.Background(), embeddableText)

	if embError == nil {

		knowledge.HasEmbedding = true

		db.UpdateOneCustom(CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": bson.M{"hasEmbedding": true}})

		COUKnowledgeEmbedding(user, knowledge.Uid, embedding)
	} else {

		db.DeleteEntity(CT_KNOWLEDGE_EMBEDDING, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		})

		db.UpdateOneCustom(CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": bson.M{"hasEmbedding": false}})

	}
}

func COUKnowledgeStep(user *models.User, knowledge *Knowledge) {

}

func DeleteKnowledgeEmbedding(user *models.User, uid primitive.ObjectID) {
	db.DeleteEntities(CT_KNOWLEDGE_EMBEDDING, bson.M{
		"_id": uid,
		"org": user.Org,
	})
}

func COUKnowledgeEmbedding(user *models.User, knUid primitive.ObjectID, embedding llmModels.Embedding) {
	DeleteKnowledgeEmbedding(user, knUid)

	kne := KnowledgeEmbedding{}
	kne.Embedding = embedding
	kne.Uid = knUid
	kne.Org = user.Org

	db.InsertEntity(CT_KNOWLEDGE_EMBEDDING, kne)
}

func GetKnowledge(user *models.User, knowledgeUid primitive.ObjectID) *Knowledge {
	return db.QueryEntity[Knowledge](CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}

func DeleteKnowledge(user *models.User, knowledgeUid primitive.ObjectID) {

	DeleteKnowledgeEmbedding(user, knowledgeUid)

	db.DeleteEntity(CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}
