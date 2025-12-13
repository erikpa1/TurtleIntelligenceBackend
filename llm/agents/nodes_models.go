package agents

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhaseType uint8

const (
	AGENT_PHASE_TRIGGER PhaseType = 0
	AGENT_PHASE_CONTROL PhaseType = 1
	AGENT_PHASE_OUTPUT  PhaseType = 2
)

type LLMAgentNode struct {
	Uid         primitive.ObjectID            `json:"uid" bson:"_id,omitempty"`
	Name        string                        `json:"name"`
	Parent      primitive.ObjectID            `json:"parent"`
	Org         primitive.ObjectID            `json:"org"`
	PosX        float32                       `json:"posX" bson:"posX"`
	PosY        float32                       `json:"posY" bson:"posY"`
	Type        string                        `json:"type"`
	PhaseType   PhaseType                     `json:"phaseType" bson:"phaseType"`
	TypeData    bson.M                        `json:"typeData" bson:"typeData"`
	Connections map[string]primitive.ObjectID `json:"connections"` //Connections are deleted by editor and it has to modify nodes
}

type LLMAgentConnection struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Parent       primitive.ObjectID `json:"parent"`
	Source       primitive.ObjectID `json:"source"`
	SourceHandle string             `json:"sourceHandle" bson:"sourceHandle"`
	Target       primitive.ObjectID `json:"target"`
	TargetHandle string             `json:"targetHandle" bson:"targetHandle"`
	Priority     int8               `json:"priority"`
	Org          primitive.ObjectID `json:"org"`
}
