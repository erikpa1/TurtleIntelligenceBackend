package items

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListItems(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListItems(user))
}

func _GetItem(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetItem(user, uid))
}

func _COUItem(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	item := tools.ObjFromJsonPtr[Item](c.PostForm("data"))
	COUItem(user, item)
	serverKit.ReturnOkJson(c, item)
}

func _DeleteItem(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteItem(user, uid)
}

func InitItemsApi(r *gin.Engine) {
	r.GET("/api/inventory/items", auth.LoginRequired, _ListItems)
	r.GET("/api/inventory/item", auth.LoginRequired, _GetItem)

	r.POST("/api/inventory/item", auth.LoginRequired, _COUItem)
	r.DELETE("/api/inventory/item", auth.LoginRequired, _DeleteItem)
}
