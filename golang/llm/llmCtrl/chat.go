package llmCtrl

import (
	"go.mongodb.org/mongo-driver/bson"
	"turtle/db"
	"turtle/llm/llmModels"
)

const CT_LLM_CHAT_HISTORY = "llm_chat_history"

func QueryChatHistory(query bson.M) []*llmModels.ChatHistory {
	return db.QueryEntities[llmModels.ChatHistory](CT_LLM_CHAT_HISTORY, query)
}
