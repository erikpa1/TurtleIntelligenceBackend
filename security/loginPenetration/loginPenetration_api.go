package loginPenetration

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _RunPenetrationTesting(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	RunLoginPenTest(user, uid)
}

func InitLoginPenetration(r *gin.Engine) {

	r.POST("/api/security/login-pentesting", auth.LoginOrApp, _RunPenetrationTesting)

	dynamicModules.InitDefaultEntitiesApi(r, "security", "login-penetration")
}
