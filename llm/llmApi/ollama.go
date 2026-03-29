package llmApi

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/llm/llmCtrl"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListOllama(c *gin.Context) {
	serverKit.ReturnOkJson(c, llmCtrl.OllamaList())
}

func _InstallOllama(c *gin.Context) {
	cluster, _ := tools.StringToObjectID(c.PostForm("cluster"))
	model := c.PostForm("model")
	serverKit.ReturnOkJson(c, llmCtrl.InstallOllama(cluster, model))
}

func InitOllamaApi(r *gin.Engine) {
	r.GET("/api/ollama/list", auth.LoginRequired, _ListOllama)
	r.POST("/api/ollama/install", auth.LoginRequired, _InstallOllama)
}
