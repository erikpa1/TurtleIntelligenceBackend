package agentTools

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/tools"
)

func _ListAgentTools(c *gin.Context) {
	tools.AutoReturn(c, ListAgentTools())

}

func InitAgentToolsApi(r *gin.Engine) {

	r.GET("/api/agent-tools", auth.AdminRequired, _ListAgentTools)

}
