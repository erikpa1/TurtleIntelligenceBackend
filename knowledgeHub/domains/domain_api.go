package domains

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/tools"
)

func _ListDomains(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListDomains(user))
}

func _COUDomain(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[Domain](c.PostForm("data"))
	COUDomain(user, data)
}

func _DeleteDomain(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	domain := tools.MongoObjectIdFromQuery(c)
	DeleteDomain(user, domain)
}

func InitDomainApi(r *gin.Engine) {
	r.GET("/api/kh/domains", auth.LoginOrApp, _ListDomains)
	r.POST("/api/kh/domain", auth.LoginOrApp, _COUDomain)
	r.DELETE("/api/kh/domain", auth.LoginOrApp, _DeleteDomain)
}
