package apiApp

import (
	"github.com/gin-gonic/gin"
	"turtle/auth"
	"turtle/models"
)

func _AddWorldEntity(c *gin.Context) {

	entity := models.NewEntity()
	entity.Name = "Gorm entity"

}

func _SaveWorld(c *gin.Context) {
	//TODO
}

func init_api_world(r *gin.Engine) {
	r.GET("/api/w/add", auth.LoginRequired, _AddWorldEntity)
	r.GET("/api/w/save", auth.LoginRequired, _SaveWorld)
}
