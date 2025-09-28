package forecasting

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForecastType int8

const (
	FORECAST_TYPE_NEURAL_NETWORK ForecastType = 0
	FORECAST_TYPE_EXP_SMOOTHING  ForecastType = 1
	FORECAST_TYPE_REGRESSION     ForecastType = 2
	FORECAST_TYPE_WINTERS        ForecastType = 3
	FORECAST_TYPE_DELPHI         ForecastType = 4
)

type Forecast struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name"`
	Type ForecastType       `json:"type"`
}

type ForecastingMethods struct {
	Type    ForecastType `json:"type"`
	Name    string       `json:"name"`
	Enabled bool         `json:"enabled"`
}
