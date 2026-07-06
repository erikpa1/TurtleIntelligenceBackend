package demand

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListDemand(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListDemand(user))
}

func _GetDemand(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetDemand(user, uid))
}

func _COUDemand(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	d := tools.ObjFromJsonPtr[DemandOrder](c.PostForm("data"))
	COUDemand(user, d)
	serverKit.ReturnOkJson(c, d)
}

func _DeleteDemand(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteDemand(user, uid)
}

func InitDemandApi(r *gin.Engine) {
	r.GET("/api/manufacturing/demands", auth.LoginRequired, _ListDemand)
	r.GET("/api/manufacturing/demand", auth.LoginRequired, _GetDemand)

	r.POST("/api/manufacturing/demand", auth.LoginRequired, _COUDemand)
	r.DELETE("/api/manufacturing/demand", auth.LoginRequired, _DeleteDemand)
}
