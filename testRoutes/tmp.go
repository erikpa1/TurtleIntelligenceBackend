package testRoutes

import (
	"turtle/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _Here(c *gin.Context) {
	lgr.Error("Here")
	tools.AutoReturn(c, bson.M{"hello": "world"})
}

func InitTestRoutes(r *gin.Engine) {
	r.GET("/api/fischertechnik", _Here)
}
