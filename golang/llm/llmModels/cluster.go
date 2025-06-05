package llmModels

import "go.mongodb.org/mongo-driver/bson/primitive"

type LLMCluster struct {
	Uid     primitive.ObjectID `json:"uid" bson:"_id"`
	Name    string             `json:"name"`
	Url     string             `json:"url"`
	xApiKey string             `json:"xApiKey"`
	Org     primitive.ObjectID `json:"org"`
}
