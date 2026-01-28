package knowledge

import (
	"context"
	"turtle/core/users"

	"turtle/db"
	"turtle/knowledgeHub/cts"
	"turtle/knowledgeHub/node"
	"turtle/lgr"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListKnowledge(user *users.User) []*Knowledge {
	return QueryKnowledge(user, bson.M{"org": user.Org})
}

func QueryKnowledge(user *users.User, query bson.M) []*Knowledge {

	findOptions := options.FindOptions{}
	findOptions.Projection = bson.M{"typeData": 0}

	query["org"] = user.Org

	lgr.InfoJson(query)

	return db.QueryEntities[Knowledge](cts.CT_KNOWLEDGE, query, &findOptions)
}

func COUKnowledge(user *users.User, knowledge *Knowledge) {
	knowledge.Org = user.Org

	if knowledge.Uid.IsZero() {
		knowledge.Uid = primitive.NewObjectID()

		db.InsertEntity(cts.CT_KNOWLEDGE, knowledge)
	} else {
		db.UpdateOneCustom(cts.CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": knowledge})
	}

	knowledgeTypeData := KnowledgePlainTextTypeData(knowledge.TypeData)
	embeddableText := knowledge.Name + " " + knowledge.Description + " " + knowledgeTypeData.GetEmbeddableString()
	embedding, embError := llmCtrl.CreateStringEmbedding(context.Background(), embeddableText)

	if embError == nil {

		knowledge.HasEmbedding = true

		db.UpdateOneCustom(cts.CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": bson.M{"hasEmbedding": true}})

		COUKnowledgeEmbedding(user, knowledge.Uid, embedding)
	} else {

		db.DeleteEntity(cts.CT_KNOWLEDGE_EMBEDDINGS, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		})

		db.UpdateOneCustom(cts.CT_KNOWLEDGE, bson.M{
			"_id": knowledge.Uid,
			"org": knowledge.Org,
		}, bson.M{"$set": bson.M{"hasEmbedding": false}})

	}
}

func COUKnowledgeStep(user *users.User, knowledge *Knowledge) {

}

func DeleteKnowledgeEmbedding(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntities(cts.CT_KNOWLEDGE_EMBEDDINGS, bson.M{
		"_id": uid,
		"org": user.Org,
	})
}

func COUKnowledgeEmbedding(user *users.User, knUid primitive.ObjectID, embedding llmModels.Embedding) {
	DeleteKnowledgeEmbedding(user, knUid)

	kne := KnowledgeEmbedding{}
	kne.Embedding = embedding
	kne.Uid = knUid
	kne.Org = user.Org

	db.InsertEntity(cts.CT_KNOWLEDGE_EMBEDDINGS, kne)
}

func GetKnowledge(user *users.User, knowledgeUid primitive.ObjectID) *Knowledge {
	return db.QueryEntity[Knowledge](cts.CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}

func DeleteKnowledge(user *users.User, knowledgeUid primitive.ObjectID) {

	DeleteKnowledgeEmbedding(user, knowledgeUid)

	db.DeleteEntity(cts.CT_KNOWLEDGE, bson.M{
		"_id": knowledgeUid,
		"org": user.Org,
	})
}

func DeleteKnowledgeOfDomain(user *users.User, domainUid primitive.ObjectID) {
	node.DeleteNodesOfDomain(user, domainUid)

	db.DeleteEntities(cts.CT_KNOWLEDGE_EMBEDDINGS, bson.M{
		"domain": domainUid,
	})

	db.DeleteEntities(cts.CT_KNOWLEDGE, bson.M{
		"domain": domainUid,
	})
}
