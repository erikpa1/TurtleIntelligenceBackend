package api

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _ListScenes(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	parent := tools.MongoObjectIdFromQueryByKey(c, "parent")
	tools.AutoReturn(c, ctrl.ListScenes(user.Org, parent))
}

func _GetScene(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, ctrl.GetScene(user.Org, uid))
}

func _COUScene(c *gin.Context) {
	user := auth.GetUserFromContext(c)

	scene := tools.ObjFromJsonPtr[models.TurtleScene](c.PostForm("data"))

	ctrl.COUScene(user, scene)
	tools.AutoReturn(c, scene.Uid)

}

func initApiScenes(r *gin.Engine) {
	r.GET("/api/scenes", auth.LoginOrApp, _ListScenes)
	r.GET("/api/scene", auth.LoginOrApp, _GetScene)
	r.POST("/api/scene", auth.LoginRequired, _COUScene)
}
