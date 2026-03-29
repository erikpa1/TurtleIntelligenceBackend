package testRoutes

import (
	"turtle/core/lgr"
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _Here(c *gin.Context) {
	lgr.Error("Here")
	serverKit.ReturnOkJson(c, bson.M{"hello": "world"})
}

func InitTestRoutes(r *gin.Engine) {
	r.GET("/api/fischertechnik", _Here)
}
