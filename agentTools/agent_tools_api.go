package agentTools

import (
	"turtle/auth"
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
)

func _ListAgentTools(c *gin.Context) {
	serverKit.ReturnOkJson(c, ListAgentsTools())

}

func InitAgentToolsApi(r *gin.Engine) {

	r.GET("/api/agent-tools", auth.AdminRequired, _ListAgentTools)

}
