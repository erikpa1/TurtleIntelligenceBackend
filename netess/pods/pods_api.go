package pods

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _ListPods(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListPods(user))
}

func _GetPod(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnJsonOr404(c, GetPod(user, uid))
}

func _COUPod(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[NetessPod](c.PostForm("data"))
	COUPod(user, data)
}

func _DeletePod(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeletePod(user, uid)
}

func InitPodsApi(r *gin.Engine) {
	r.GET("/api/netess/pods", auth.LoginOrApp, _ListPods)
	r.GET("/api/netess/pod", auth.LoginOrApp, _GetPod)
	r.POST("/api/netess/pod", auth.LoginOrApp, _COUPod)
	r.DELETE("/api/netess/pod", auth.LoginOrApp, _DeletePod)

}
