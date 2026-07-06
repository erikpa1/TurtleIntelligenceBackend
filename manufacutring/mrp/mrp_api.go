package mrp

import (
	"turtle/auth"
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
)

func _RunMrp(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, RunMrp(user))
}

func InitMrpApi(r *gin.Engine) {
	// The run is a pure, stateless computation over the current master data,
	// exposed on both verbs for convenience.
	r.POST("/api/manufacturing/mrp/run", auth.LoginRequired, _RunMrp)
	r.GET("/api/manufacturing/mrp/run", auth.LoginRequired, _RunMrp)
}
