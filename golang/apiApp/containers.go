package apiApp

import (
	"github.com/erikpa1/turtle/ctrlApp"
	"github.com/erikpa1/turtle/modelsApp"
	"github.com/erikpa1/turtle/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
