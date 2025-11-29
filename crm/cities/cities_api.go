package cities

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

func InitCitiesApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "city")
}
