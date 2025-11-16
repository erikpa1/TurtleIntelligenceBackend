package domains

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org         primitive.ObjectID `json:"org"`
	Name        string             `json:"name"`
	Color       string             `json:"color"`
	Icon        string             `json:"icon"`
	Description string             `json:"description"`
	CreatedAt   tools.Milliseconds `json:"createdAt" bson:"createdAt"`
}
