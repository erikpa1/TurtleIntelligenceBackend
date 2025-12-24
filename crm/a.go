package crm

import (
	"github.com/gin-gonic/gin"
	"turtle/crm/cities"
	"turtle/crm/contracts"
	"turtle/crm/invoices"
)

func InitCrmApi(r *gin.Engine) {
	cities.InitCitiesApi(r)
	invoices.InitInvoicesApi(r)
	contracts.InitContractsApi(r)

}
