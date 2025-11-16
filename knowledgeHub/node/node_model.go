package node

import "go.mongodb.org/mongo-driver/bson/primitive"

type KnowledgeModel struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Artifact    string             `json:"artifact"`
}
