package llmCtrl

import (
	"fmt"
	"github.com/erikpa1/turtle/credentials"
	"github.com/erikpa1/turtle/db"
	"github.com/erikpa1/turtle/llm/llmModels"
	"github.com/erikpa1/turtle/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_LLM_CLUSTERS = "llm_clusters"

func ListLLMClusters(org primitive.ObjectID) []*llmModels.LLMCluster {
	return db.QueryEntities[llmModels.LLMCluster](CT_LLM_CLUSTERS, bson.M{
		"org": org,
	})
}

func DeleteLLMCluster(user *models.User, uid primitive.ObjectID) {

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

func COULLMCluster(user *models.User, cluster *llmModels.LLMCluster) {

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

func GetLLMCluster(user *models.User, uid primitive.ObjectID) *llmModels.LLMCluster {

	if uid.IsZero() {
		return &llmModels.LLMCluster{Url: fmt.Sprintf("localhost:%s", credentials.GetPort())}
	} else {
		return db.QueryEntity[llmModels.LLMCluster](CT_LLM_CLUSTERS, bson.M{
			"org": uid,
			"_id": user.Org,
		})
	}

}
