package llmCtrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"turtle/db"
	"turtle/llm/llmModels"
	"turtle/models"
)

const CT_LLM_MODELS = "llm_models"

var ROUND_ROBIN_CLUSTERS = map[primitive.ObjectID]int{}
var ROUND_ROBIN_CLUSTERS_MUT = sync.Mutex{}

func ListLLMModels(org primitive.ObjectID) []*llmModels.LLM {
	return db.QueryEntities[llmModels.LLM](CT_LLM_MODELS, bson.M{
		"org": org,
	})
}
func GetLLMModel(user *models.User, uid primitive.ObjectID) *llmModels.LLM {
	return db.QueryEntity[llmModels.LLM](CT_LLM_MODELS,
		bson.M{
			"org": user.Org,
			"_id": uid,
		})
}

func DeleteModelsOfCluster(user *models.User, uid primitive.ObjectID) {
	if user.IsAdmin() {
		db.DeleteEntities(CT_LLM_MODELS, bson.M{
			"org": user.Org,
			"_id": uid,
		})
	}
}

func DeleteLLMModel(user *models.User, uid primitive.ObjectID) {
	if user.IsAdminWithError() {
		db.DeleteEntities(CT_LLM_MODELS, bson.M{
			"org": user.Org,
			"_id": uid,
		})
	}
}

func COULLMModel(user *models.User, model *llmModels.LLM) {
	if user.IsAdmin() {
		if model.Uid.IsZero() {
			model.Org = user.Org
			db.InsertEntity(CT_LLM_MODELS, model)
		} else {
			db.UpdateOneCustom(CT_LLM_MODELS,
				bson.M{
					"org": user.Org,
					"_id": model.Uid,
				},
				bson.M{"$set": model},
			)
		}

		if model.IsDefault {
			db.UpdateEntitiesWhere(CT_LLM_MODELS,
				bson.M{"_id": bson.M{
					"$nin": bson.A{model.Uid},
				}},
				bson.M{"isDefault": false},
			)
		}

	}
}

func ListAgenticOrNormalModels() []*llmModels.LLM {
	agents := db.QueryEntities[llmModels.LLM](CT_LLM_MODELS, bson.M{"isAgentic": true})

	if len(agents) > 0 {
		return agents
	} else {
		return db.QueryEntities[llmModels.LLM](CT_LLM_MODELS, bson.M{"isAgentic": false})
	}
}

func GetRoundRobinCluster(clusters []primitive.ObjectID, modelUid primitive.ObjectID) primitive.ObjectID {

	if len(clusters) > 0 {
		ROUND_ROBIN_CLUSTERS_MUT.Lock()

		tmp, exists := ROUND_ROBIN_CLUSTERS[modelUid]

		if exists {
			if tmp >= len(clusters) {
				ROUND_ROBIN_CLUSTERS[modelUid] = 0
				return clusters[0]
			} else {
				ROUND_ROBIN_CLUSTERS[modelUid] += 1
				return clusters[tmp]
			}
		} else {
			ROUND_ROBIN_CLUSTERS[modelUid] = 1
			return clusters[0]
		}

		ROUND_ROBIN_CLUSTERS_MUT.Unlock()
	}

	return primitive.ObjectID{}

}
