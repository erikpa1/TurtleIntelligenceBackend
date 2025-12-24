package cities

import (
	"github.com/gin-gonic/gin"
	"turtle/dynamicModules"
)

func InitCitiesApi(r *gin.Engine) {
	dynamicModules.InitDefaultEntitiesApi(r, "crm", "city")
}
