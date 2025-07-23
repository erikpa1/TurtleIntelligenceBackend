package fn

import "go.mongodb.org/mongo-driver/bson/primitive"

type FnLight struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name"`
}

type Fn struct {
	Uid       primitive.ObjectID     `json:"uid" bson:"_id,omitempty"`
	Name      string                 `json:"name"`
	Inputs    FnInputs               `json:"inputs"`
	Output    FnOutputs              `json:"output"`
	Variables map[string]interface{} `json:"variables"`
}

type FnInputs struct {
}

type FnOutputs struct {
}

type FnVariable struct {
	Id int32 `json:"id"`
}

type FnNode struct {
}
