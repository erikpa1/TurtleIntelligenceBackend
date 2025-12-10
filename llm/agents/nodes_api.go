package agents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func _COUNodes(c *gin.Context) {

	user := auth.GetUserFromContext(c)

	type _Request struct {
		Modified []*LLMAgentNode      `json:"modified"`
		Created  []*LLMAgentNode      `json:"created"`
		Deleted  []primitive.ObjectID `json:"deleted"`
	}

	req := tools.ObjFromJsonPtr[_Request](c.PostForm("data"))

	lg.LogEson(req)

	for _, deleted := range req.Deleted {
		DeleteAgentNode(deleted)
	}

	for _, modified := range req.Modified {
		COUNode(user, modified)
	}

	InsertNodes(user, req.Created)
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
	r.POST("/api/llm/agent-nodes", auth.LoginRequired, _COUNodes)

	r.DELETE("/api/llm/agent-node", auth.LoginRequired, _DeleteNode)

	r.POST("/api/llm/agent/exec/:agentUid", auth.LoginRequired, _ExecNode)
	r.POST("/api/llm/agent/exec/:agentUid/:orgUid", auth.LoginRequired, _ExecOrgNode)
}
