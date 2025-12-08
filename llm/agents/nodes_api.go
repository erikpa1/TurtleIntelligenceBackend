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

func _ExecNode(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	agentUid, _ := tools.StringToObjectID(c.Param("agentUid"))
	ExeNodeWithUid(user.Org, agentUid)
}

func _ExecOrgNode(c *gin.Context) {
	orgUid, _ := tools.StringToObjectID(c.Param("orgUid"))
	agentUid, _ := tools.StringToObjectID(c.Param("agentUid"))
	ExeNodeWithUid(orgUid, agentUid)

}

func InitLLMAgentNodes(r *gin.Engine) {
	r.GET("/api/llm/agent-nodes/query", auth.LoginRequired, _QueryAgentNodes)
	r.POST("/api/llm/agent-node", auth.LoginRequired, _COUNode)

	r.DELETE("/api/llm/agent-node", auth.LoginRequired, _DeleteNode)

	r.POST("/api/llm/agent/exec/:agentUid", auth.LoginRequired, _ExecNode)
	r.POST("/api/llm/agent/exec/:orgUid/:agentUid", auth.LoginRequired, _ExecOrgNode)
}
