package llmModels

import "go.mongodb.org/mongo-driver/bson/primitive"

type LlmModel struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org          primitive.ObjectID `json:"uid"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	ModelVersion string             `json:"modelVersion" bson:"modelVersion"`
}
