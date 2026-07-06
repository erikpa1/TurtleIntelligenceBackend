package workcenters

import "go.mongodb.org/mongo-driver/bson/primitive"

// WorkCenter is a capacity providing resource (machine, line, cell) used by the
// APS scheduler. Routing operations are executed on work centers and the
// scheduler treats each work center as a single finite-capacity resource.
type WorkCenter struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org         primitive.ObjectID `json:"org" bson:"org"`
	Code        string             `json:"code" bson:"code"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`

	CapacityHoursPerDay float64 `json:"capacityHoursPerDay" bson:"capacityHoursPerDay"`
	// Efficiency in percent (100 = nominal). Processing time is divided by this
	// factor when scheduling.
	Efficiency  float64 `json:"efficiency" bson:"efficiency"`
	CostPerHour float64 `json:"costPerHour" bson:"costPerHour"`

	Active bool `json:"active" bson:"active"`
}
