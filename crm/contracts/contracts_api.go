package contracts

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

func InitContractsApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "contracts")
}
