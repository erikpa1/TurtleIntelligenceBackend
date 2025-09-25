package forecasting

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForecastType int8

const (
	FORECAST_TYPE_LINEAR_SMOOTHING      ForecastType = 0
	FORECAST_TYPE_EXPONENTIAL_SMOOTHING ForecastType = 1
	FORECAST_TYPE_WINTERS               ForecastType = 2
)

type Forecast struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name"`
	Type ForecastType       `json:"type"`
}
