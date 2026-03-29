package apiApp

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/core/serverKit"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/tools"
)

func _QueryActors(c *gin.Context) {
	query := tools.QueryHeader[bson.M](c)
	serverKit.ReturnOkJson(c, ctrlApp.QueryActors(query))
}
func _GetActor(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnOkJson(c, ctrlApp.GetActor(uid))
}

func _CouActor(c *gin.Context) {
	data := tools.ObjFromJsonPtr[modelsApp.Actor](c.PostForm("data"))
	ctrlApp.COUActor(data)
}

func _DeleteActor(c *gin.Context) {
	uid := serverKit.MongoObjectIdFromQuery(c)

	ctrlApp.DeleteActor(uid)
}

func _InitApiActors(r *gin.Engine) {
	r.GET("/api/actors/query", _QueryActors)
	r.GET("/api/actor", _GetActor)
	r.POST("/api/actor", _CouActor)
	r.DELETE("/api/actor", _DeleteActor)
}
