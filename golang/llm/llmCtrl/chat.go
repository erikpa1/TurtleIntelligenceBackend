package llmCtrl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strings"
	"turtle/db"
	"turtle/lg"
	"turtle/llm/llmModels"
	"turtle/models"
	"turtle/tools"
)

const CT_LLM_CHAT_HISTORY = "llm_chat_history"

func GetChat(uid primitive.ObjectID) *llmModels.ChatHistory {

	return db.QueryEntity[llmModels.ChatHistory](CT_LLM_CHAT_HISTORY, bson.M{"_id": uid})
}

func QueryChatsHistory(query bson.M) []*llmModels.ChatHistoryLight {

	findOptions := options.Find()
	findOptions.Projection = llmModels.ChatHistoryLightProjection()

	return db.QueryEntities[llmModels.ChatHistoryLight](CT_LLM_CHAT_HISTORY, query, findOptions)
}

func StartLLMChat(user *models.User, suggestedName string) primitive.ObjectID {
	tmp := llmModels.ChatHistory{}
	tmp.Uid = primitive.NewObjectID()
	tmp.At = tools.GetTimeNowMillis()
	tmp.UserUid = user.Uid
	tmp.Org = user.Org
	tmp.Name = suggestedName
	tmp.Conversation = make([]llmModels.ConversationSegment, 0)
	tmp.Answered = false

	db.InsertEntity(CT_LLM_CHAT_HISTORY, tmp)

	return tmp.Uid
}

func AddUserQuestion(user *models.User, chatId primitive.ObjectID, text string) {
	segment := llmModels.ConversationSegment{}
	segment.At = tools.GetTimeNowMillis()
	segment.IsUser = true
	segment.Text = text // Don't forget to set the text!

	db.UpdateOneCustom(CT_LLM_CHAT_HISTORY, bson.M{
		"_id":     chatId,
		"org":     user.Org,
		"userUid": user.Uid,
	},
		bson.M{
			"$push": bson.M{"conversation": segment},
			"$set":  bson.M{"answered": false},
		}, // Specify the array field
	)
}

func AddChatAnswer(user *models.User, chatId primitive.ObjectID, text string) {
	segment := llmModels.ConversationSegment{}
	segment.At = tools.GetTimeNowMillis()
	segment.IsUser = false
	segment.Text = text // Don't forget to set the text!

	segment.SmartTexts = FindAllContent(text)

	db.UpdateOneCustom(CT_LLM_CHAT_HISTORY, bson.M{
		"_id":     chatId,
		"org":     user.Org,
		"userUid": user.Uid,
	},
		bson.M{
			"$push": bson.M{"conversation": segment},
			"$set":  bson.M{"answered": true},
		}, // Specify the array field
	)
}

func DeleteChat(user primitive.ObjectID, chatUid primitive.ObjectID) {
	db.DeleteEntity(CT_LLM_CHAT_HISTORY, bson.M{"_id": chatUid, "userUid": user})
}

func AskModel(c *gin.Context, user *models.User, modelUid primitive.ObjectID, prompt string) string {

	model := GetLLMModel(user, modelUid)

	if model != nil {
		if model.Cluster.IsZero() {

			lg.LogOk(fmt.Sprintf("Going to ask localhost LLM"))

			return AskLangChainModel(c, model, prompt)
		} else {
			cluster := GetLLMCluster(user, model.Cluster)

			if cluster != nil {
				if strings.Contains(cluster.Url, "localhost") ||
					strings.Contains(cluster.Url, "127.0.0.1") ||
					strings.Contains(cluster.Url, "0.0.0.0") {

					lg.LogOk(fmt.Sprintf("Going to ask on url [%s]", cluster.Url))

					return AskLangChainModel(c, model, prompt)
				}
			} else {
				resp, err := http.Get(fmt.Sprintf("%s%s", cluster.Url, "/api/llm/ask"))

				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					lg.LogE(resp)
				}

				lg.LogE("Cluster is invalid")
			}
		}
	} else {
		lg.LogE("Model don't exists anymore")
	}

	return "--unanswered--"
}
