package api

import (
	"github.com/gin-gonic/gin"
)

func InitApi(r *gin.Engine) {
	initApiTools(r)
}
