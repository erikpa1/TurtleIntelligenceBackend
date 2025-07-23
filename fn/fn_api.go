package fn

import "github.com/gin-gonic/gin"

func _ListFunctions(c *gin.Context) {

}

func _GetFunctions(c *gin.Context) {

}

func _CreateFunction(c *gin.Context) {
}

func _DeleteFunction(c *gin.Context) {
}

func InitFnApi(r *gin.Engine) {
	r.GET("/api/functions", _ListFunctions)
	r.GET("/api/function", _GetFunctions)
	r.POST("/api/function", _CreateFunction)
	r.DELETE("/api/function", _DeleteFunction)

}
