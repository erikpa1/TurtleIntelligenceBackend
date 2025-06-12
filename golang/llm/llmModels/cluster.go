package llmModels

import "go.mongodb.org/mongo-driver/bson/primitive"

type LLMCluster struct {
	Uid     primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Url     string             `json:"url"`
	xApiKey primitive.ObjectID `json:"xApiKey"`
	Org     primitive.ObjectID `json:"org"`
}
