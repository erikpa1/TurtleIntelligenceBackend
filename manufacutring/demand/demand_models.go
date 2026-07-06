package demand

import "go.mongodb.org/mongo-driver/bson/primitive"

// DemandOrder is a line of independent demand (a sales order or a forecast
// entry) for a finished product. It is the primary input that drives the MRP
// run: gross requirements for the top level items.
type DemandOrder struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID `json:"org" bson:"org"`
	Reference string             `json:"reference" bson:"reference"`

	ProductUid string  `json:"productUid" bson:"productUid"`
	ProductSku string  `json:"productSku,omitempty" bson:"productSku,omitempty"`
	Quantity   float64 `json:"quantity" bson:"quantity"`
	Uom        string  `json:"uom,omitempty" bson:"uom,omitempty"`

	// DueDate is the date the demand must be satisfied (ISO yyyy-mm-dd).
	DueDate string `json:"dueDate" bson:"dueDate"`
	// DemandType is "sales" or "forecast".
	DemandType string `json:"demandType,omitempty" bson:"demandType,omitempty"`
	// Status is "open", "released" or "closed". Closed demand is ignored by MRP.
	Status string `json:"status,omitempty" bson:"status,omitempty"`
}
