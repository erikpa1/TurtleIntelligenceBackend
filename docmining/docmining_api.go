package docmining

import (
	"turtle/lg"

	"github.com/gin-gonic/gin"
)

func _FileToProcessPosted(c *gin.Context) {
	lg.LogE("Here")

}

func InitDocMiningApi(r *gin.Engine) {

	r.POST("/api/docmining/document", _FileToProcessPosted)
}
