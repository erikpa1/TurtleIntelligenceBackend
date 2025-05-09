package apiApp

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/ctrlApp"
	"turtle/modelsApp"
	"turtle/tools"
)

func _QueryContainers(c *gin.Context) {
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, ctrlApp.QueryContainers(query))
}

func _CreateContainer(c *gin.Context) {
	data := tools.ObjFromJsonPtr[modelsApp.Container](c.PostForm("data"))
	ctrlApp.CreateContainer(data)
}

func _InitApiContainers(r *gin.Engine) {
	r.POST("/api/container", _CreateContainer)
	r.GET("/api/containers/query", _QueryContainers)
}
