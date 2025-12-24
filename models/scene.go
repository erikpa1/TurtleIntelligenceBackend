package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
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
