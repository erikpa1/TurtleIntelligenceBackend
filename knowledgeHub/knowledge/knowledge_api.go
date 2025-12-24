package knowledge

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/tools"
)

func _QueryKnowledge(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryBsonHeader(c)
	tools.AutoReturn(c, QueryKnowledge(user, query))
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

func _COUKnowledgeStep(c *gin.Context) {
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
	r.GET("/api/kh/knowledge/query", auth.LoginOrApp, _QueryKnowledge)
	r.GET("/api/kh/knowledge", auth.LoginOrApp, _GetKnowledge)
	r.POST("/api/kh/knowledge", auth.LoginOrApp, _COUKnowledge)
	r.POST("/api/kh/knowledge/step", auth.LoginOrApp, _COUKnowledgeStep)
	r.DELETE("/api/kh/knowledge", auth.LoginOrApp, _DeleteKnowledge)
}
