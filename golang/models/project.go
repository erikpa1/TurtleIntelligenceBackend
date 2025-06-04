package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type TurtleProject struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	CreatedBy primitive.ObjectID
	Org       primitive.ObjectID

	Name        string
	Description string

	CreatedAt tools.Milliseconds
	UpdatedAt tools.Milliseconds
}
