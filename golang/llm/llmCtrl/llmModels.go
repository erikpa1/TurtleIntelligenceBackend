package llmCtrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/llm/llmModels"
	"turtle/models"
)

const CT_LLM_MODELS = "llm_models"

func ListLLMModels(org primitive.ObjectID) []*llmModels.LlmModel {
	return db.QueryEntities[llmModels.LlmModel](CT_LLM_MODELS, bson.M{
		"org": org,
	})
}
func GetLLMModel(user *models.User, uid primitive.ObjectID) *llmModels.LlmModel {
	return db.QueryEntity[llmModels.LlmModel](CT_LLM_MODELS,
		bson.M{
			"org": user.Org,
			"uid": uid,
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
	if user.IsAdmin() {
		db.DeleteEntities(CT_LLM_MODELS, bson.M{
			"org": user.Org,
			"_id": uid,
		})
	}
}

func COULLMModel(user *models.User, model *llmModels.LlmModel) {
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
	}
}
