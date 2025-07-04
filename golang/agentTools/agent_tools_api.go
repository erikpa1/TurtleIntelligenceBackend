package agentTools

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/tools"
)

func _ListAgentTools(c *gin.Context) {
	tools.AutoReturn(c, bson.M{})

}

func InitAgentToolsApi(r *gin.Engine) {

	r.GET("/api/agent-tools", auth.AdminRequired, _ListAgentTools)

}
