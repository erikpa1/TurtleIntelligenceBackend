package modelsApp

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type Document struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Extension    string             `json:"extension"`
	CreatedAt    tools.Milliseconds `json:"createdAt" bson:"createdAt"`
	UpdatedAt    tools.Milliseconds `json:"updatedAt" bson:"updatedAt"`
	HasEmbedding bool               `json:"hasEmbedding" bson:"hasEmbedding"`
}

type DocumentEmbedding struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Embedding []float64          `json:"embeding"`
}
