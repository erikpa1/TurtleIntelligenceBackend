package forecasting

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/tables"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _QueryForecasts(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := tools.QueryHeader[bson.M](c)
	tools.AutoReturn(c, QueryForecasts(user, query))
}

func _GetForecast(c *gin.Context) {

}

func _COUForecast(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	obj := tools.ObjFromJsonPtr[Forecast](c.PostForm("data"))
	COUForecast(user, obj)
}

func _DeleteForecast(c *gin.Context) {
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteForecast(uid)
}

func _ListForecastingMethods(c *gin.Context) {
	tools.AutoReturn(c, ListForecastingMethods())
}

func InitForecastingApi(r *gin.Engine) {
	r.GET("/api/forecasts/query", auth.LoginRequired, _QueryForecasts)
	r.GET("/api/forecast", auth.LoginRequired, _GetForecast)
	r.GET("/api/forecast/methods", auth.LoginRequired, _ListForecastingMethods)
	r.POST("/api/forecast", auth.LoginRequired, _COUForecast)
	r.DELETE("/api/forecast", auth.LoginRequired, _DeleteForecast)

	tables.CreateGinRouting(r, "forecasting")

}
