package routing

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListRoutings(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListRoutings(user))
}

func _GetRouting(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetRouting(user, uid))
}

func _COURouting(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	rt := tools.ObjFromJsonPtr[Routing](c.PostForm("data"))
	COURouting(user, rt)
	serverKit.ReturnOkJson(c, rt)
}

func _DeleteRouting(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteRouting(user, uid)
}

func InitRoutingApi(r *gin.Engine) {
	r.GET("/api/manufacturing/routings", auth.LoginRequired, _ListRoutings)
	r.GET("/api/manufacturing/routing", auth.LoginRequired, _GetRouting)

	r.POST("/api/manufacturing/routing", auth.LoginRequired, _COURouting)
	r.DELETE("/api/manufacturing/routing", auth.LoginRequired, _DeleteRouting)
}
