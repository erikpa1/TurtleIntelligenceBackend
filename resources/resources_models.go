package resources

import "go.mongodb.org/mongo-driver/bson/primitive"

type Resource struct {
	Uid  primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name string             `json:"name"`
	Org  primitive.ObjectID `json:"org"`
	Unit string             `json:"unit"`
}
