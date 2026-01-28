package llmCtrl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"turtle/core/users"
	"turtle/db"
	"turtle/lgr"
	"turtle/llm/llmModels"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func StartLLMChat(user *users.User, suggestedName string) primitive.ObjectID {
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

func AddUserQuestion(user *users.User, chatId primitive.ObjectID, text string) {
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

func AddChatAnswer(user *users.User, chatId primitive.ObjectID, text string) {
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

func AskModelStream(c *gin.Context, user *users.User, modelUid primitive.ObjectID, prompt string) {
	model := GetLLMModel(user, modelUid)

	if model != nil {
		if len(model.Clusters) == 0 {
			lgr.Ok(fmt.Sprintf("Going to ask on local model [%s]", model.ModelVersion))

			AskLangChainModelStream(c, model, prompt)
		} else {

			clusterUid := GetRoundRobinCluster(model.Clusters, model.Uid)

			cluster := GetLLMCluster(user, clusterUid)

			if cluster != nil {
				if strings.Contains(cluster.Url, "localhost") ||
					strings.Contains(cluster.Url, "127.0.0.1") ||
					strings.Contains(cluster.Url, "0.0.0.0") {

					lgr.Ok(fmt.Sprintf("Going to ask on url [%s]", cluster.Url))

					AskLangChainModelStream(c, model, prompt)
				}
			} else {
				resp, err := http.Get(fmt.Sprintf("%s%s", cluster.Url, "/api/llm/ask"))

				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					lgr.ErrorJson(resp)
				}

				lgr.Error("Cluster is invalid")
			}
		}
	} else {
		lgr.Error("Model don't exists anymore")
	}

}

func AskModelForDescription(c *gin.Context,
	user *users.User,
	modelUid primitive.ObjectID,
	userQuery string,
	maxWords int,
) string {
	finalPrompt := fmt.Sprintf(`
SYSTEM: You are document analyst, you describe document with %d words maximum
INSTRUCTIONS:
1. Analyze the user query
2. Extract the text
3. Respond in the following JSON format:
{
  "description": "description",
}
USER QUERY: {%s}
`, userQuery, maxWords)

	response := AskModel(c, user, modelUid, finalPrompt)

	maybeJson, _, err := tools.FindFirstJsonString(response)

	if err != nil {
		lgr.Error(err.Error())
	}

	converted := bson.M{}

	err = json.Unmarshal([]byte(maybeJson), &converted)

	if err != nil {
		lgr.Error(err.Error())
		return ""
	} else {
		description, ok := converted["description"].(string)

		if ok {
			return description
		} else {
			lgr.Error("No description found")
			return ""
		}
	}

}

func AskModel(c *gin.Context, user *users.User, modelUid primitive.ObjectID, prompt string) string {

	model := GetLLMOrDefault(user, modelUid)

	if model != nil {
		if len(model.Clusters) == 0 {
			lgr.Ok(fmt.Sprintf("Going to ask localhost LLM: %s", model.ModelVersion))
			return AskLangChainModel(c, model, prompt)
		} else {
			return AskModelRemote(c, user, model, prompt)
		}
	} else {
		lgr.Error("Model don't exists anymore [", modelUid, "]")
		tools.AutoNotFound(c, "llm.notfound")
	}

	return "--unanswered--"
}

func AskModelRemote(c *gin.Context, user *users.User, model *llmModels.LLM, prompt string) string {
	clusterUid := GetRoundRobinCluster(model.Clusters, model.Uid)

	cluster := GetLLMCluster(user, clusterUid)

	if cluster != nil {
		if strings.Contains(cluster.Url, "localhost") ||
			strings.Contains(cluster.Url, "127.0.0.1") ||
			strings.Contains(cluster.Url, "0.0.0.0") {

			lgr.Ok(fmt.Sprintf("Going to ask on url [%s]", cluster.Url))

			return AskLangChainModel(c, model, prompt)
		}
	} else {
		resp, err := http.Get(fmt.Sprintf("%s%s", cluster.Url, "/api/llm/ask"))

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			lgr.ErrorJson(resp)
		}

		lgr.Error("Cluster is invalid")
	}

	return ""
}
