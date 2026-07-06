package workcenters

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListWorkCenters(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListWorkCenters(user))
}

func _GetWorkCenter(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetWorkCenter(user, uid))
}

func _COUWorkCenter(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	w := tools.ObjFromJsonPtr[WorkCenter](c.PostForm("data"))
	COUWorkCenter(user, w)
	serverKit.ReturnOkJson(c, w)
}

func _DeleteWorkCenter(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteWorkCenter(user, uid)
}

func InitWorkCentersApi(r *gin.Engine) {
	r.GET("/api/manufacturing/workcenters", auth.LoginRequired, _ListWorkCenters)
	r.GET("/api/manufacturing/workcenter", auth.LoginRequired, _GetWorkCenter)

	r.POST("/api/manufacturing/workcenter", auth.LoginRequired, _COUWorkCenter)
	r.DELETE("/api/manufacturing/workcenter", auth.LoginRequired, _DeleteWorkCenter)
}
