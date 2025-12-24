package llmCtrl

import (
	"fmt"
	"turtle/core/users"
	"turtle/credentials"
	"turtle/db"
	"turtle/llm/llmModels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_LLM_CLUSTERS = "llm_clusters"

func ListLLMClusters(org primitive.ObjectID) []*llmModels.LLMCluster {
	return db.QueryEntities[llmModels.LLMCluster](CT_LLM_CLUSTERS, bson.M{
		"org": org,
	})
}

func DeleteLLMCluster(user *users.User, uid primitive.ObjectID) {

	if user.IsAdmin() {
		DeleteModelsOfCluster(user, uid)

		db.DeleteEntity(CT_LLM_CLUSTERS,
			bson.M{
				"org": user.Org,
				"_id": uid,
			},
		)
	}
}

func COULLMCluster(user *users.User, cluster *llmModels.LLMCluster) {

	if user.IsAdminWithError() {

		cluster.Org = user.Org

		if cluster.Uid.IsZero() {
			db.InsertEntity(CT_LLM_CLUSTERS, cluster)
		} else {
			db.UpdateOneCustom(CT_LLM_CLUSTERS, bson.M{
				"org": user.Org,
				"_id": cluster.Uid,
			},
				bson.M{"$set": cluster},
			)
		}
	}
}

func GetLLMCluster(user *users.User, uid primitive.ObjectID) *llmModels.LLMCluster {

	if uid.IsZero() {
		return &llmModels.LLMCluster{Url: fmt.Sprintf("localhost:%s", credentials.GetPort())}
	} else {
		return db.QueryEntity[llmModels.LLMCluster](CT_LLM_CLUSTERS, bson.M{
			"org": uid,
			"_id": user.Org,
		})
	}

}
