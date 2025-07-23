package nn

import "github.com/gin-gonic/gin"

func _COUNN(c *gin.Context) {

}

func _ListModels(c *gin.Context) {

}

func _DeleteNN(c *gin.Context) {

}

func InitNNApi(r *gin.Engine) {
	r.POST("/api/sim-model", _COUNN)
	r.GET("/api/models", _ListModels)
	r.DELETE("/api/model", _DeleteNN)
}
