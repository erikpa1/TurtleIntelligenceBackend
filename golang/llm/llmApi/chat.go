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
	llm, err := ollama.New(ollama.WithModel("deepseek-coder-v2:latest"))

	if err == nil {
		completion, complErr := llms.GenerateFromSinglePrompt(c, llm, "how are you today?")

		if complErr == nil {
			lg.LogI(completion)
		}

	} else {
		lg.LogE(err)
	}

}

func InitLLMChatApi(r *gin.Engine) {
	r.GET("/api/llm/chat-history", auth.LoginRequired, _GetChatHistory)
	r.POST("/api/llm/chat-ask", auth.LoginRequired, _ChatAsk)
}
