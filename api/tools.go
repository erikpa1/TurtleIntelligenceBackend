package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"turtle/auth"
)

func _GetMongoId(c *gin.Context) {
	c.String(http.StatusOK, primitive.NewObjectID().Hex())
}

func initApiTools(r *gin.Engine) {
	r.GET("/api/tools/mongoid", auth.LoginOrApp, _GetMongoId)

}
