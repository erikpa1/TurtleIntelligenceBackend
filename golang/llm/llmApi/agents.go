package llmApi

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"
)

func _ListLLMAgents(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.ListLLMAgents(user))
}

func _DeleteLLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	llmCtrl.DeleteLLMAgent(user, tools.MongoObjectIdFromQuery(c))
}

func _COULLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	agent := tools.ObjFromJson[llmModels.LLMAgent](c.PostForm("data"))
	llmCtrl.COULLMAgent(user, &agent)
}

func _TestLLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	text := c.PostForm("text")
	agent, _ := primitive.ObjectIDFromHex(c.PostForm("agent"))

	tools.AutoReturn(c, llmCtrl.TestLLMAgent(c, user, agent, text))

}

func InitLLMAgents(r *gin.Engine) {
	r.GET("/api/llm/agents", auth.LoginRequired, _ListLLMAgents)

	r.GET("/api/llm/agent", auth.LoginRequired, _TestLLMAgent)

	r.POST("/api/llm/agent", auth.LoginRequired, _COULLMAgent)
	r.POST("/api/llm/agent/test", auth.LoginRequired, _TestLLMAgent)

	r.DELETE("/api/llm/agent", auth.LoginRequired, _DeleteLLMAgent)

}
