package llmApi

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/llm/llmCtrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _ListOllama(c *gin.Context) {
	tools.AutoReturn(c, llmCtrl.OllamaList())
}

func _InstallOllama(c *gin.Context) {
	cluster, _ := tools.StringToObjectID(c.PostForm("cluster"))
	model := c.PostForm("model")
	tools.AutoReturn(c, llmCtrl.InstallOllama(cluster, model))
}

func InitOllamaApi(r *gin.Engine) {
	r.GET("/api/ollama/list", auth.LoginRequired, _ListOllama)
	r.POST("/api/ollama/install", auth.LoginRequired, _InstallOllama)
}
