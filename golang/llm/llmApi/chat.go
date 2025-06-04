package llmApi

import (
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/auth"
	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/tools"
)

func _GetChat(c *gin.Context) {
	objectUid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, llmCtrl.GetChat(objectUid))
}

func _ChatAsk(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	conversation, _ := primitive.ObjectIDFromHex(c.PostForm("chatUid"))
	text := c.PostForm("text")

	llmCtrl.AddUserQuestion(user, conversation, text)

	ollmodel := ollama.WithModel("deepseek-coder-v2:latest")
	keepAlive := ollama.WithKeepAlive("10h")

	llm, err := ollama.New(ollmodel, keepAlive)

	if err == nil {
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, text)
		if complErr == nil {
			lg.LogI(completion)

			llmCtrl.AddChatAnswer(user, conversation, completion)
		} else {
			lg.LogE(completion)
		}
	} else {
		lg.LogE(err)
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

func InitLLMChatApi(r *gin.Engine) {
	r.GET("/api/llm/chat", auth.LoginRequired, _GetChat)
	r.GET("/api/llm/chats", auth.LoginRequired, _GetChatsHistory)
	r.POST("/api/llm/chat-ask", auth.LoginRequired, _ChatAsk)
	r.POST("/api/llm/chat/start", auth.LoginRequired, _StartChat)

	r.DELETE("/api/llm/chat", auth.LoginRequired, _DeleteChat)

}
