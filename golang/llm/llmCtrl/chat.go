package llmCtrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/db"
	"turtle/llm/llmModels"
	"turtle/tools"
)

const CT_LLM_CHAT_HISTORY = "llm_chat_history"

func QueryChatHistory(query bson.M) []*llmModels.ChatHistory {
	return db.QueryEntities[llmModels.ChatHistory](CT_LLM_CHAT_HISTORY, query)
}

func StartLLMChat() primitive.ObjectID {
	tmp := llmModels.ChatHistory{}
	tmp.Uid = primitive.NewObjectID()
	tmp.At = tools.Milliseconds(tools.GetTimeNowMillis())

	db.InsertEntity(CT_LLM_CHAT_HISTORY, tmp)

	return tmp.Uid
}
