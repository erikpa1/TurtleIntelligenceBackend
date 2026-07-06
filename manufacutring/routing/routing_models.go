package routing

import "go.mongodb.org/mongo-driver/bson/primitive"

// Operation is a single step in a routing, executed on a work center. Total
// processing time for a production order of quantity Q is:
//
//	SetupMinutes + RunMinutesPerUnit * Q
type Operation struct {
	Sequence          int     `json:"sequence" bson:"sequence"`
	Name              string  `json:"name" bson:"name"`
	WorkCenterUid     string  `json:"workCenterUid" bson:"workCenterUid"`
	WorkCenterName    string  `json:"workCenterName,omitempty" bson:"workCenterName,omitempty"`
	SetupMinutes      float64 `json:"setupMinutes" bson:"setupMinutes"`
	RunMinutesPerUnit float64 `json:"runMinutesPerUnit" bson:"runMinutesPerUnit"`
}

// Routing describes the ordered operations required to manufacture a product.
// It is consumed by the APS scheduler together with the planned production
// orders produced by MRP.
type Routing struct {
	Uid        primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org        primitive.ObjectID `json:"org" bson:"org"`
	Code       string             `json:"code" bson:"code"`
	Name       string             `json:"name" bson:"name"`
	ProductUid string             `json:"productUid" bson:"productUid"`
	ProductSku string             `json:"productSku,omitempty" bson:"productSku,omitempty"`
	Operations []Operation        `json:"operations" bson:"operations"`
}
