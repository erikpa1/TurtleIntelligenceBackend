package invoices

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

func InitInvoicesApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "invoices")
}
