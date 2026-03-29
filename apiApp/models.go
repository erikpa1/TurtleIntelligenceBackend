package apiApp

import (
	"turtle/auth"
	"turtle/core/serverKit"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/tools"

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
	serverKit.ReturnOkJson(c, ctrlApp.ListModels())
}

func _InitApiModels(r *gin.Engine) {
	r.POST("/api/sims/models", auth.LoginRequired, _CouModel)
	r.GET("/api/sims/models", auth.LoginRequired, _ListModels)
	r.DELETE("/api/sims/models", auth.LoginRequired, _DeleteModel)

}
