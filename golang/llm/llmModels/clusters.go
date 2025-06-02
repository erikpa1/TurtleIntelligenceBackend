package llmModels

import "go.mongodb.org/mongo-driver/bson/primitive"

type LLMCluster struct {
	Uid primitive.ObjectID `bson:"_id" json:"uid"`
}
