package apiApp

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/auth"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/tools"
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
	r.POST("/api/sim-models", auth.LoginRequired, _CouModel)
	r.GET("/api/sim-models", auth.LoginRequired, _ListModels)
	r.DELETE("/api/sim-models", auth.LoginRequired, _DeleteModel)

}
