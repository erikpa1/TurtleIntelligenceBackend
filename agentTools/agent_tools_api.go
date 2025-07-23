package agentTools

import (
	"github.com/erikpa1/turtle/auth"
	"github.com/erikpa1/turtle/tools"
	"github.com/gin-gonic/gin"
)

func _ListAgentTools(c *gin.Context) {
	tools.AutoReturn(c, ListAgentsTools())

}

func InitAgentToolsApi(r *gin.Engine) {

	r.GET("/api/agent-tools", auth.AdminRequired, _ListAgentTools)

}
