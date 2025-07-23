package llmApi

import "github.com/gin-gonic/gin"

func InitLLMApi(r *gin.Engine) {
	InitLLMChatApi(r)
	InitLlmAndClusters(r)
	InitLLMAgents(r)
	InitOllamaApi(r)
}
