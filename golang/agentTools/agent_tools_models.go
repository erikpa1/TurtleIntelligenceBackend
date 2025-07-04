package agentTools

import "go.mongodb.org/mongo-driver/bson/primitive"

type AgentTool struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
	Inputs      string             `json:"inputType" bson:"inputType"`
}
