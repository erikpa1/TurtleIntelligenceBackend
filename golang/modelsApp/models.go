package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

type Model struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
