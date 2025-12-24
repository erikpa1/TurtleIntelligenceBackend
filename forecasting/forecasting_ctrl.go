package forecasting

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_FORECASTS = "forecasts"
const CT_FORECASTS_RESULTS = "forecasts_results"
const CT_FORECASTS_DATA = "forecasts_data"

func QueryForecasts(user *users.User, query bson.M) []*Forecast {
	return db.QueryEntities[Forecast](CT_FORECASTS, user.FillOrgQuery(query))
}

func COUForecast(user *users.User, forecast *Forecast) {

	forecast.Org = user.Org

	if forecast.Uid.IsZero() {
		db.InsertEntity(CT_FORECASTS, forecast)
	} else {
		db.SetByOrgAndId(CT_FORECASTS, forecast.Uid, user.Org, forecast)
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
