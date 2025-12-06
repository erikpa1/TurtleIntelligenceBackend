package agents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _QueryAgentNodes(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryBsonHeader(c)
	tools.AutoReturn(c, QueryNodes(user, query))
}

func _COUNode(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[LLMAgentNode](c.PostForm("data"))
	COUNode(user, data)
}

func _DeleteNode(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteAgentNode(uid)
}

func InitLLMAgentNodes(r *gin.Engine) {
	r.GET("/api/llm/agent-nodes/query", auth.LoginRequired, _QueryAgentNodes)
	r.POST("/api/llm/agent-node", auth.LoginRequired, _COUNode)
	r.DELETE("/api/llm/agent-node", auth.LoginRequired, _DeleteNode)
}
