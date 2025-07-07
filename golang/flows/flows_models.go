package flows

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type FlowExecution struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org       primitive.ObjectID `json:"org"`
	Flow      primitive.ObjectID `json:"flow"`
	At        tools.Milliseconds `json:"at"`
	Status    int8               `json:"status"`
	Callstack string             `json:"callStack"`
}

type Flow struct {
	Uid    primitive.ObjectID     `json:"uid" bson:"_id,omitempty"`
	Org    primitive.ObjectID     `json:"org" bson:"org"`
	Name   string                 `json:"name"`
	States map[string]interface{} `json:"states"`
}
