package apiApp

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrlApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _QueryActors(c *gin.Context) {
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, ctrlApp.QueryActors(query))
}
func _GetActor(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	tools.AutoReturn(c, ctrlApp.GetActor(uid))
}

func _CouActor(c *gin.Context) {
	data := tools.ObjFromJsonPtr[modelsApp.Actor](c.PostForm("data"))
	ctrlApp.COUActor(data)
}

func _DeleteActor(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)

	ctrlApp.DeleteActor(uid)
}

func _InitApiActors(r *gin.Engine) {
	r.GET("/api/actors/query", _QueryActors)
	r.GET("/api/actor", _GetActor)
	r.POST("/api/actor", _CouActor)
	r.DELETE("/api/actor", _DeleteActor)
}
