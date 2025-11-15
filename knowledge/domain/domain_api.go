package domain

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/gin-gonic/gin"
)

func _ListDomains(c *gin.Context) {

}

func _COUDomain(c *gin.Context) {

}

func _DeleteDomain(c *gin.Context) {

}

func InitDomainApi(r *gin.Engine) {
	r.GET("/api/domains", auth.LoginOrApp, _ListDomains)
	r.POST("/api/domain", auth.LoginOrApp, _COUDomain)
	r.DELETE("/api/domain", auth.LoginOrApp, _DeleteDomain)
}
