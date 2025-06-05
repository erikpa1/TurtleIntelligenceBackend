package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type XApiKey struct {
	Uid   primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name  string             `json:"name"`
	Value string             `json:"value"`
}
