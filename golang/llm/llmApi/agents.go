package llmApi

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
)

func InitLLMAgents(r *gin.Engine) {
	r.GET("/api/llm/clusters", auth.LoginRequired, _ListLLMClusters)

}
