package bisSubjects

import (
	"github.com/gin-gonic/gin"
	"turtle/dynamicModules"
)

func InitBisSubjectsApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "subjects")
}
