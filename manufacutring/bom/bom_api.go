package bom

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListBoms(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListBoms(user))
}

func _GetBom(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetBom(user, uid))
}

func _COUBom(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	b := tools.ObjFromJsonPtr[Bom](c.PostForm("data"))
	COUBom(user, b)
	serverKit.ReturnOkJson(c, b)
}

func _DeleteBom(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteBom(user, uid)
}

func InitBomApi(r *gin.Engine) {
	r.GET("/api/manufacturing/boms", auth.LoginRequired, _ListBoms)
	r.GET("/api/manufacturing/bom", auth.LoginRequired, _GetBom)

	r.POST("/api/manufacturing/bom", auth.LoginRequired, _COUBom)
	r.DELETE("/api/manufacturing/bom", auth.LoginRequired, _DeleteBom)
}
