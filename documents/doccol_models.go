package documents

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentsCollection struct {
	Uid          primitive.ObjectID    `json:"uid" bson:"_id,omitempty"`
	Org          primitive.ObjectID    `json:"org" bson:"org"`
	Name         string                `json:"name"`
	Filter       string                `json:"filter"`
	SelectedTags []string              `json:"selectedTags"`
	Items        []*primitive.ObjectID `json:"items"`
}
