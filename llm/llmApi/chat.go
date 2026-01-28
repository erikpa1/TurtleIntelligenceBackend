package llmApi

import (
	"turtle/auth"
	"turtle/lgr"
	"turtle/llm/llmCtrl"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _GetChat(c *gin.Context) {
	objectUid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, llmCtrl.GetChat(objectUid))
}

func _AskModel(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	modelUid, _ := primitive.ObjectIDFromHex(c.PostForm("modelUid"))
	text := c.PostForm("text")
	tools.AutoReturn(c, llmCtrl.AskModel(c, user, modelUid, text))
}

func _AskModelStream(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	modelUid, _ := primitive.ObjectIDFromHex(c.Query("modelUid"))
	text := c.Query("text")
	llmCtrl.AskModelStream(c, user, modelUid, text)

}
func _ChatAsk(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	modelUid, _ := primitive.ObjectIDFromHex(c.PostForm("modelUid"))
	conversation, _ := primitive.ObjectIDFromHex(c.PostForm("chatUid"))
	text := c.PostForm("text")
	isAgent := c.Query("isAgentChat") == "true"

	llmCtrl.AddUserQuestion(user, conversation, text)

	if isAgent {
		respone := llmCtrl.ChatAgent(c, user, modelUid, text)
		lgr.ErrorJson(respone)
	} else {
		completion := llmCtrl.AskModel(c, user, modelUid, text)
		llmCtrl.AddChatAnswer(user, conversation, completion)
	}

}

func _StartChat(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	name := c.PostForm("name")
	tools.AutoReturn(c, llmCtrl.StartLLMChat(user, name))

}

func _GetChatsHistory(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.QueryChatsHistory(bson.M{
		"userUid": user.Uid,
	}))
}

func _DeleteChat(c *gin.Context) {
	chatUid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)
	llmCtrl.DeleteChat(user.Uid, chatUid)
}

func _TestEmbeding(c *gin.Context) {
	llmCtrl.ExampleEmbedding()
}

func InitLLMChatApi(r *gin.Engine) {
	r.GET("/api/llm/chat", auth.LoginRequired, _GetChat)
	r.GET("/api/llm/chats", auth.LoginRequired, _GetChatsHistory)
	r.POST("/api/llm/chat-ask", auth.LoginRequired, _ChatAsk) //
	r.POST("/api/llm/ask", auth.LoginRequired, _AskModel)     // "/api/llm/ask"

	r.POST("/api/llm/chat/start", auth.LoginRequired, _StartChat)
	r.POST("/api/llm/embedding", auth.LoginRequired, _TestEmbeding)

	r.GET("/api/llm/ask/stream", auth.LoginRequired, _AskModelStream) // "/api/llm/ask"

	r.DELETE("/api/llm/chat", auth.LoginRequired, _DeleteChat)

}
