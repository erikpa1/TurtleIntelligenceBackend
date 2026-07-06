package inventory

import (
	"turtle/inventory/items"

	"github.com/gin-gonic/gin"
)

// InitInventoryApi wires all inventory sub-modules (materials/items, ...).
func InitInventoryApi(r *gin.Engine) {
	items.InitItemsApi(r)
}
