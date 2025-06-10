package api

import (
	"github.com/gin-gonic/gin"
)

func InitApi(r *gin.Engine) {
	initUsersApi(r)
	initApiProjects(r)
	initApiScenes(r)
	initApiTools(r)
	InitDocumentsApi(r)
}
