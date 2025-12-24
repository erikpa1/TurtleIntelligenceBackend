package tags

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/tools"
)

func _ListTags(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, ListTags(user, query))
}

func _DeleteTag(c *gin.Context) {
	uid := c.Query("uid")
	user := auth.GetUserFromContext(c)
	DeleteTag(user, uid)
}

func _CreateTag(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tag := tools.ObjFromJsonPtr[Tag](c.PostForm("data"))
	CreateTag(user, tag)
}

func _UpdateTag(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tag := tools.ObjFromJsonPtr[Tag](c.PostForm("data"))
	UpdateTag(user, tag)
}

func InitTagsApi(c *gin.Engine) {
	c.GET("/api/tags", auth.LoginOrApp, _ListTags)
	c.DELETE("/api/tags", auth.LoginOrApp, _DeleteTag)
	c.PUT("/api/tag", auth.LoginOrApp, _UpdateTag)
	c.POST("/api/tag", auth.LoginOrApp, _CreateTag)
}
