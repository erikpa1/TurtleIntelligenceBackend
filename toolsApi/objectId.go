package toolsApi

import (
	"net/http"
	"turtle/auth"
	"turtle/lg"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _TranslateString(c *gin.Context) {
	uid := c.Query("uid")
	converted, _ := tools.StringToObjectID(uid)
	hexResult := converted.Hex()
	lg.LogE(hexResult)
	c.String(http.StatusOK, hexResult)

}

func InitObjectIdApi(r *gin.Engine) {
	r.GET("/api/object-id/from-str", auth.LoginOrApp, _TranslateString)
}
