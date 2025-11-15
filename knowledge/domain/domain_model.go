package domain

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CreatedAt   tools.Milliseconds `json:"createdAt" bson:"createdAt"`
}
