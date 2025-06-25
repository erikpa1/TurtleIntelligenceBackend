package knowledge

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/tools"
)

func _ListKnowledge(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListKnowledge(user))
}

func _GetKnowledge(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	knowUid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, GetKnowledge(user, knowUid))

}

func _COUKnowledge(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	dataStr := c.PostForm("data")

	data := tools.ObjFromJson[Knowledge](dataStr)
	COUKnowledge(user, &data)
}

func _DeleteKnowledge(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	knowUid := tools.MongoObjectIdFromQuery(c)
	DeleteKnowledge(user, knowUid)

}

func InitKnowledgeApi(r *gin.Engine) {

	r.GET("/api/knowledge/list", auth.LoginOrApp, _ListKnowledge)
	r.GET("/api/knowledge", auth.LoginOrApp, _GetKnowledge)
	r.POST("/api/knowledge", auth.LoginOrApp, _COUKnowledge)
	r.DELETE("/api/knowledge", auth.LoginOrApp, _DeleteKnowledge)
}
