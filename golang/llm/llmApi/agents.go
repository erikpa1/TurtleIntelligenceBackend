package llmApi

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/tools"
)

func _ListLLMAgents(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.ListLLMAgents(user))
}

func _TestLLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	text := c.PostForm("text")

	llmCtrl.TestLLMAgent(c, user, text)
}

func InitLLMAgents(r *gin.Engine) {
	r.GET("/api/llm/agents", auth.LoginRequired, _ListLLMAgents)
	r.GET("/api/llm/agents/test", auth.LoginRequired, _TestLLMAgent)

}
