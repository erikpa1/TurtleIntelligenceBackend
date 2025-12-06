package crm

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/crm/cities"
	"github.com/erikpa1/TurtleIntelligenceBackend/crm/contracts"
	"github.com/erikpa1/TurtleIntelligenceBackend/crm/invoices"
	"github.com/gin-gonic/gin"
)

func InitCrmApi(r *gin.Engine) {
	cities.InitCitiesApi(r)
	invoices.InitInvoicesApi(r)
	contracts.InitContractsApi(r)

}
