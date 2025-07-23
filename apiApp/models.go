package apiApp

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/ctrlApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/modelsApp"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func _CouModel(c *gin.Context) {
	data := tools.VecFromJStr[modelsApp.Model](c.PostForm("data"))

	for _, model := range data {
		ctrlApp.COUModel(&model)
	}
}

func _DeleteModel(c *gin.Context) {
	uid, _ := primitive.ObjectIDFromHex(c.Query("uid"))
	ctrlApp.DeleteModel(uid)
}

func _ListModels(c *gin.Context) {
	tools.AutoReturn(c, ctrlApp.ListModels())
}

func _InitApiModels(r *gin.Engine) {
	r.POST("/api/sim-models", _CouModel)
	r.GET("/api/sim-models", _ListModels)
	r.DELETE("/api/sim-models", _DeleteModel)

}
