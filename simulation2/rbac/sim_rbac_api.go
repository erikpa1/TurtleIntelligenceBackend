package rbac

import (
	"turtle/auth"
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetRbacGraphBuilder() {

}

func _ListRbacRoles(c *gin.Context) {
	serverKit.ReturnOkJson(c, bson.A{
		RBAC_SIM_READER,
		RBAC_SIM_EDITOR,
		RBAC_SIM_ADMIN,
	})
}

func InitSimRbacRolesApi(r *gin.Engine) {
	r.GET("/api/sims/rbac", auth.LoginRequired, _ListRbacRoles)
}
