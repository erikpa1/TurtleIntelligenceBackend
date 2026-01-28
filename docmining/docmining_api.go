package docmining

import (
	"turtle/lgr"

	"github.com/gin-gonic/gin"
)

func _FileToProcessPosted(c *gin.Context) {
	lgr.Error("Here")

}

func InitDocMiningApi(r *gin.Engine) {

	r.POST("/api/docmining/document", _FileToProcessPosted)
}
