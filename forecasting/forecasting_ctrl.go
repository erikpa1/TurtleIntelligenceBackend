package forecasting

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_FORECASTS = "forecasts"
const CT_FORECASTS_RESULTS = "forecasts_results"
const CT_FORECASTS_DATA = "forecasts_data"

func QueryForecasts(user *models.User, query bson.M) []*Forecast {
	return db.QueryEntities[Forecast](CT_FORECASTS, user.FillOrgQuery(query))
}

func COUForecast(forecast *Forecast) {
	if forecast.Uid.IsZero() {
		db.InsertEntity(CT_FORECASTS, forecast)
	} else {
		db.SetById(CT_FORECASTS, forecast.Uid, forecast)
	}
}

func DeleteForecast(uid primitive.ObjectID) {
	db.DeleteEntitiesOfParent(CT_FORECASTS_DATA, uid)
	db.DeleteEntitiesOfParent(CT_FORECASTS_RESULTS, uid)
	db.DeleteEntityWithUid(CT_FORECASTS, uid)
}

func ListForecastingMethods() []ForecastingMethods {
	return []ForecastingMethods{
		{
			Type:    FORECAST_TYPE_EXP_SMOOTHING,
			Name:    "exp.smoothing",
			Enabled: true,
		},
		{
			Type:    FORECAST_TYPE_NEURAL_NETWORK,
			Name:    "nn",
			Enabled: true,
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
