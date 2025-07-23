package flows

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _ListFlows(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListFlows(user))

}

func _GetFlow(c *gin.Context) {
	flowUid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, GetFlow(user, flowUid))
}

func _DeleteFlow(c *gin.Context) {
	flowUid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)
	DeleteFlow(user, flowUid)
}

func _COUFlow(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[Flow](c.PostForm("data"))
	COUFlow(user, data)
}

func _CallFlow(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	status, errStr := CallFlow(user, tools.MongoObjectIdFromQuery(c))

	tools.AutoReturn(c, bson.M{
		"status":  status,
		"respnse": errStr,
	})

}

func InitFlowsApi(r *gin.Engine) {
	r.GET("/api/flows", auth.LoginOrApp, _ListFlows)
	r.GET("/api/flow", auth.LoginOrApp, _GetFlow)
	r.POST("/api/flow", auth.LoginOrApp, _COUFlow)
	r.DELETE("/api/flow", auth.LoginOrApp, _DeleteFlow)

	r.Any("/api/flow/call/:id", auth.LoginOrApp, _CallFlow)
}
