package apiApp

import "github.com/gin-gonic/gin"

func InitApiApp(r *gin.Engine) {
	init_api_world(r)
	_InitApiContainers(r)
	_InitApiModels(r)
	_InitApiActors(r)
}
