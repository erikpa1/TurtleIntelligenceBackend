package loginPenetration

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/dynamicModules"
	"turtle/tools"
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
