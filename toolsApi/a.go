package toolsApi

import "github.com/gin-gonic/gin"

func InitToolsApi(r *gin.Engine) {
	InitObjectIdApi(r)
}
