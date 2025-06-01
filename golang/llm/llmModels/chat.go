package llmModels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type ChatHistory struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	At   tools.Milliseconds `json:"at" bson:"at"`
}
