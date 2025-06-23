package llmModels

import "go.mongodb.org/mongo-driver/bson/primitive"

type LLM struct {
	Uid                        primitive.ObjectID   `json:"uid" bson:"_id,omitempty"`
	Org                        primitive.ObjectID   `json:"org"`
	Clusters                   []primitive.ObjectID `json:"clusters"`
	IsDefault                  bool                 `json:"isDefault" bson:"isDefault"`
	Preload                    bool                 `json:"preload" bson:"preload"`
	Name                       string               `json:"name"`
	Description                string               `json:"description"`
	ModelVersion               string               `json:"modelVersion" bson:"modelVersion"`
	Ttl                        string               `json:"ttl"`
	IsAgentic                  bool                 `json:"isAgentic" bson:"isAgentic"`
	DefaultTemperature         float32              `json:"defaultTemperature" bson:"defaultTemperature"`
	CanUserOverrideTemperature bool                 `json:"canUserOverrideTemperature" bson:"canUserOverrideTemperature"`
}
