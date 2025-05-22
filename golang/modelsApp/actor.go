package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

type Actor struct {
	Uid   primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Model primitive.ObjectID `json:"model" bson:"model"`
	Color string             `json:"color" bson:"color"`
}
