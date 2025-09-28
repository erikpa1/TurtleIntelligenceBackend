package forecasting

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
)

const CT_FORECASTS = "forecasts"
const CT_FORECASTS_RESULTS = "forecasts_results"
const CT_FORECASTS_DATA = "forecasts_data"

func QueryForecasts(user *models.User, query bson.M) []*Forecast {
	return db.QueryEntities[Forecast](CT_FORECASTS, user.FillOrgQuery(query))
}

func ListForecastingMethods() []ForecastingMethods {
	return []ForecastingMethods{
		{
			Type:    FORECAST_TYPE_NEURAL_NETWORK,
			Name:    "nn",
			Enabled: true,
		},
		{
			Type: FORECAST_TYPE_EXP_SMOOTHING,
			Name: "exp.smoothing",
		},
		{
			Type: FORECAST_TYPE_REGRESSION,
			Name: "regression",
		},
		{
			Type: FORECAST_TYPE_WINTERS,
			Name: "Winters-Holt",
		},
		{
			Type: FORECAST_TYPE_DELPHI,
			Name: "Delphi",
		},
	}
}
