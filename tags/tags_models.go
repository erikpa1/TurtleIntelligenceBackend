package tags

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid         string             `json:"uid"`
	Name        string             `json:"name"`
	Org         primitive.ObjectID `json:"org"`
	Color       string             `json:"color"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
}
