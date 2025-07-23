package models

import (
	"github.com/erikpa1/turtle/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TurtleScene struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID
	Parent    primitive.ObjectID
	CreatedBy primitive.ObjectID
	UpdatedBy primitive.ObjectID
	CreatedAt tools.Milliseconds
	UpdatedAt tools.Milliseconds

	Name        string
	Description string
	Type        string
}
