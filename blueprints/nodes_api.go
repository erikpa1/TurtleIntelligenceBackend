package blueprints

import (
	"turtle/auth"
	"turtle/blueprints/models"
	"turtle/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _QueryAgentNodes(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryBsonHeader(c)

	tools.AutoReturn(c, QueryNodes(user, query))
}

func _QueryAgentEdges(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryBsonHeader(c)

	tools.AutoReturn(c, QueryEdges(user, query))
}

func _COUNode(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[models.Node](c.PostForm("data"))
	COUNode(user, data)
}

func _COUNodes(c *gin.Context) {

	user := auth.GetUserFromContext(c)

	type _Request struct {
		Modified     []*models.Node       `json:"modified"`
		Created      []*models.Node       `json:"created"`
		Deleted      []primitive.ObjectID `json:"deleted"`
		NewEdges     []*models.NodeEdge   `json:"newEdges"`
		DeletedEdges []primitive.ObjectID `json:"deletedEdges"`
	}

	req := tools.ObjFromJsonPtr[_Request](c.PostForm("data"))

	for _, deleted := range req.Deleted {
		DeleteNode(deleted)
	}

	for _, modified := range req.Modified {
		COUNode(user, modified)
	}

	if len(req.Created) > 0 {
		InsertNodes(user, req.Created)
	}

	if len(req.NewEdges) > 0 {
		InsertEdges(user, req.NewEdges)
	}

	if len(req.DeletedEdges) > 0 {
		DeleteEdges(user, bson.M{"_id": bson.M{"$in": req.DeletedEdges}})
	}

}

func _DeleteNode(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteNode(uid)
}

func _PlayAgentNode(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	nodeUid := tools.MongoObjectIdFromQuery(c)

	playNodeContext := models.NodePlayContext{
		Gin:         c,
		User:        user,
		IsLocalHost: c.RemoteIP() == "::1",
	}

	lgr.ErrorJson(nodeUid)

	PlayNode(&playNodeContext, nodeUid)

	tools.AutoReturn(c, bson.M{
		"pipeline": playNodeContext.Pipeline,
	})
}

func _PlayDiagram(c *gin.Context) {
	nodeUid := tools.MongoObjectIdFromQuery(c)
	PlayBlueprint(c, nodeUid)
}

func InitBlueprintsApi(r *gin.Engine) {
	InitNodesLibrary()

	r.GET("/api/blueprint/all/prompt", auth.LoginRequired, _GetAllAgentsPrompt)
	r.GET("/api/blueprints", auth.LoginRequired, _ListBlueprints)
	r.GET("/api/blueprint", auth.LoginRequired, _TestLLMAgent)
	r.POST("/api/blueprint", auth.LoginRequired, _COULLMAgent)
	r.POST("/api/blueprint/test", auth.LoginRequired, _TestLLMAgent)
	r.DELETE("/api/blueprint", auth.LoginRequired, _DeleteBlueprint)

	//Nodes
	r.GET("/api/blueprints/nodes/query", auth.LoginRequired, _QueryAgentNodes)
	r.GET("/api/blueprints/edges/query", auth.LoginRequired, _QueryAgentEdges)
	r.POST("/api/blueprints/node", auth.LoginRequired, _COUNode)

	r.POST("/api/blueprints/nodes", auth.LoginRequired, _COUNodes)

	r.DELETE("/api/blueprints/node", auth.LoginRequired, _DeleteNode)

	//Play
	r.POST("/api/blueprints/play", _PlayAgentNode) //_PlayDiagram) //_PlayAgentNode

}
