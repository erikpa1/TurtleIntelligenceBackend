package llmApi

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/tools"
)

func _ListOllama(c *gin.Context) {
	tools.AutoReturn(c, llmCtrl.OllamaList())
}

func InitOllamaApi(r *gin.Engine) {
	r.GET("/api/ollama/list", auth.LoginRequired, _ListOllama)
}
