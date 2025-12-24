package llm

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/auth"
	"turtle/tools"
)

func _ListIncidents(c *gin.Context) {

	tools.AutoReturn(c, bson.M{})

}

func InitIncidentsApi(r *gin.Engine) {

	r.GET("/api/llm/incidents", auth.AdminRequired, _ListIncidents)

}
