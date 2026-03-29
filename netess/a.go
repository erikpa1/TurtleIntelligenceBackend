package netess

import (
	"turtle/netess/pods"

	"github.com/gin-gonic/gin"
)

func InitNetessApi(r *gin.Engine) {
	pods.InitPodsApi(r)
}
