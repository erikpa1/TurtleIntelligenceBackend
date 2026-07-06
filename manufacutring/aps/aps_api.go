package aps

import (
	"turtle/auth"
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
)

func _RunAps(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, RunAps(user))
}

func InitApsApi(r *gin.Engine) {
	r.POST("/api/manufacturing/aps/run", auth.LoginRequired, _RunAps)
	r.GET("/api/manufacturing/aps/run", auth.LoginRequired, _RunAps)
}
