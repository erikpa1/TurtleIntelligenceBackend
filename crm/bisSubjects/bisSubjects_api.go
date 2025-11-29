package bisSubjects

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/dynamicModules"
	"github.com/gin-gonic/gin"
)

func InitBisSubjectsApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "subjects")
}
