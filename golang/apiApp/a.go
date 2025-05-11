package apiApp

import "github.com/gin-gonic/gin"

func InitApiApp(r *gin.Engine) {
	init_api_world(r)
	init_api_worldsim(r)
	_InitApiContainers(r)
	_InitApiModels(r)
}
