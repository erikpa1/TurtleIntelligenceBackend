package modelsApp

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entity struct {
	Uid          primitive.ObjectID            `json:"uid" bson:"_id,omitempty"`
	Model        primitive.ObjectID            `json:"model" bson:"model"`
	Dependencies map[string]primitive.ObjectID `json:"dependencies" bson:"dependencies"`
	Type         string                        `json:"type" bson:"type"`
	Name         string                        `json:"name" bson:"name"`
	Position     [3]float32                    `json:"position" bson:"position"`
	VisType      string                        `json:"visType" bson:"visType"`
	VisUid       string                        `json:"visUid" bson:"visUid"`
	TypeData     map[string]interface{}        `json:"typeData" bson:"typeData"`
}

type EntityConnection struct {
	Uid      primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Model    primitive.ObjectID `json:"model" bson:"model"`
	A        primitive.ObjectID `json:"a" bson:"a"`
	B        primitive.ObjectID `json:"b" bson:"b"`
	IsTwoWay bool               `json:"isTwoWay" bson:"isTwoWay"`
}
