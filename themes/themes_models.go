package themes

import "go.mongodb.org/mongo-driver/bson/primitive"

type ThemeLight struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	IsDefault bool               `json:"isDefault"`
}
