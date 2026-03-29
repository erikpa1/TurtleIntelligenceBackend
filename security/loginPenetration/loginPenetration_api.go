package loginPenetration

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/dynamicModules"

	"github.com/gin-gonic/gin"
)

func _RunPenetrationTesting(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	RunLoginPenTest(user, uid)
}

func InitLoginPenetration(r *gin.Engine) {

	r.POST("/api/security/login-pentesting", auth.LoginOrApp, _RunPenetrationTesting)

	dynamicModules.InitDefaultEntitiesApi(r, "security", "login-penetration")
}
