package forecasting

import "github.com/gin-gonic/gin"

func _ListForecasts(c *gin.Context) {

}

func _GetForecast(c *gin.Context) {

}

func _COUForecast(c *gin.Context) {
}

func _DeleteForecast(c *gin.Context) {
}

func InitForecastingApi(r *gin.Engine) {
	r.GET("/api/forecasts", _ListForecasts)
	r.GET("/api/forecast", _GetForecast)
	r.POST("/api/forecast", _COUForecast)
	r.DELETE("/api/forecast", _DeleteForecast)

}
