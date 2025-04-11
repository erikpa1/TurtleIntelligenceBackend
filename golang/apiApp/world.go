package apiApp

import (
	"github.com/gin-gonic/gin"
	"turtle/models"
)

func _AddWorldEntity(c *gin.Context) {

	entity := models.NewEntity()
	entity.Name = "Gorm entity"

}

func init_api_world(r *gin.Engine) {
	r.GET("/api/w/add", _AddWorldEntity)
}
