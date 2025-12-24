package invoices

import (
	"github.com/gin-gonic/gin"
	"turtle/dynamicModules"
)

func InitInvoicesApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "invoices")
}
