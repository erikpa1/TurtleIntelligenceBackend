package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

type Container struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Dimension []float32          `json:"dimension" bson:"dimension"`
	MaxHeight float32            `json:"maxHeight" bson:"maxHeight"`
}
