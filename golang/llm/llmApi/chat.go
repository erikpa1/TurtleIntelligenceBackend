package llmApi

import (
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/tools"
)

func _GetChatHistory(c *gin.Context) {
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, llmCtrl.QueryChatHistory(query))
}

func _ChatAsk(c *gin.Context) {

	text := c.PostForm("text")

	llm, err := ollama.New(ollama.WithModel("deepseek-coder-v2:latest"))

	if err == nil {
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, text)

		if complErr == nil {
			lg.LogI(completion)
		}

	} else {
		lg.LogE(err)
	}

}

func _StartChat(c *gin.Context) {
	tools.AutoReturn(c, llmCtrl.StartLLMChat())
}

func InitLLMChatApi(r *gin.Engine) {
	r.GET("/api/llm/chat-history", auth.LoginRequired, _GetChatHistory)
	r.POST("/api/llm/chat-ask", auth.LoginRequired, _ChatAsk)
	r.POST("/api/llm/chat/start", auth.LoginRequired, _StartChat)
}
