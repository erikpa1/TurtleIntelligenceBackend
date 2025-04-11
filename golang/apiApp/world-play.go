package apiApp

import "github.com/gin-gonic/gin"

func _PlayWorld(c *gin.Context) {

}

func init_api_worldsim(r *gin.Engine) {
	r.GET("/api/w/simulate", _PlayWorld)
}
