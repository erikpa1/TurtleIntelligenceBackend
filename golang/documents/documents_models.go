package documents

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/llm/llmModels"
	"turtle/tools"
)

type RightsLevel int8

const (
	RIGHTS_LEVEL_ORGANIZATION             = 0
	RIGHTS_LEVEL_GROUP        RightsLevel = 1
	RIGHTS_LEVEL_USER         RightsLevel = 2
)

type Document struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org          primitive.ObjectID `json:"org"`
	User         primitive.ObjectID `json:"user"`
	RightsLevel  RightsLevel        `json:"rightsLevel" bson:"rightsLevel"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Extension    string             `json:"extension"`
	CreatedAt    tools.Milliseconds `json:"createdAt" bson:"createdAt"`
	UpdatedAt    tools.Milliseconds `json:"updatedAt" bson:"updatedAt"`
	HasEmbedding bool               `json:"hasEmbedding" bson:"hasEmbedding"`
}

type DocumentEmbedding struct {
	Uid           primitive.ObjectID  `json:"uid" bson:"_id,omitempty"`
	Org           primitive.ObjectID  `json:"org""`
	Embedding     llmModels.Embedding `json:"embedding"`
	DescEmbedding llmModels.Embedding `json:"descEmbedding" bson:"descEmbedding"`
}

type DocumentExtraction struct {
	Uid        primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org        primitive.ObjectID `json:"org" bson:"org"`
	Extraction string             `json:"extraction"`
}
