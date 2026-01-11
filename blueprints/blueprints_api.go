package blueprints

import (
	"turtle/auth"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _ListBlueprints(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, llmCtrl.ListLLMAgents(user))
}

func _DeleteBlueprint(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	uid := tools.MongoObjectIdFromQuery(c)

	DeleteNodesOfBlueprint(user, uid)
	llmCtrl.DeleteLLMAgent(user, uid)
}

func _COULLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	agent := tools.ObjFromJsonPtr[llmModels.LLMAgent](c.PostForm("data"))
	llmCtrl.COULLMAgent(user, agent)
}

func _TestLLMAgent(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	text := c.PostForm("text")
	agent, _ := primitive.ObjectIDFromHex(c.PostForm("agent"))

	tools.AutoReturn(c, llmCtrl.ChatAgent(c, user, agent, text))

}

func _GetAllAgentsPrompt(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	tools.AutoReturn(c, llmCtrl.GetOverallAgentsPrompt(user, "[your query here]"))
}
