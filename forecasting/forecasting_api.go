package forecasting

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/tableData"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _ListForecasts(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, QueryForecasts(user, query))
}

func _GetForecast(c *gin.Context) {

}

func _COUForecast(c *gin.Context) {
}

func _DeleteForecast(c *gin.Context) {
}

func _ListForecastingMethods(c *gin.Context) {
	tools.AutoReturn(c, ListForecastingMethods())
}

func InitForecastingApi(r *gin.Engine) {
	r.GET("/api/forecasts", _ListForecasts)
	r.GET("/api/forecast", _GetForecast)
	r.GET("/api/forecast/methods", _ListForecastingMethods)
	r.POST("/api/forecast", _COUForecast)
	r.DELETE("/api/forecast", _DeleteForecast)

	tableData.CreateGinRouting(r, "forecasting")

}
