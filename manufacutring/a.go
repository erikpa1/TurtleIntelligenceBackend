package manufacturing

import (
	"turtle/manufacutring/aps"
	"turtle/manufacutring/bom"
	"turtle/manufacutring/demand"
	"turtle/manufacutring/mrp"
	"turtle/manufacutring/routing"
	"turtle/manufacutring/workcenters"

	"github.com/gin-gonic/gin"
)

// InitManufacturingApi wires all manufacturing sub-modules: the bill of
// materials, the planning master data (demand, work centers, routings) and the
// planning engines (MRP and APS).
func InitManufacturingApi(r *gin.Engine) {
	bom.InitBomApi(r)

	// Advanced planning master data.
	demand.InitDemandApi(r)
	workcenters.InitWorkCentersApi(r)
	routing.InitRoutingApi(r)

	// Planning engines.
	mrp.InitMrpApi(r)
	aps.InitApsApi(r)
}
