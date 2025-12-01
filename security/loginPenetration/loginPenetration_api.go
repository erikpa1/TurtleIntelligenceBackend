package loginPenetration

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

func InitLoginPenetration(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "security", "login-penetration")
}
