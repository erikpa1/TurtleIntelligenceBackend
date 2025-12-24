package contracts

import (
	"github.com/gin-gonic/gin"
	"turtle/dynamicModules"
)

func InitContractsApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "contracts")
}
