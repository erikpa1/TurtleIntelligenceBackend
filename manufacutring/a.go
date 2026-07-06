package manufacturing

import (
	"turtle/manufacutring/bom"

	"github.com/gin-gonic/gin"
)

// InitManufacturingApi wires all manufacturing sub-modules (BOM, ...).
func InitManufacturingApi(r *gin.Engine) {
	bom.InitBomApi(r)
}
