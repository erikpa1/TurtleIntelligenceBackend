package llmModels

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type LLMAgent struct {
	Uid            primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Org            primitive.ObjectID `json:"org" bson:"org,omitempty"`
	UserLevel      int8               `json:"userLevel"` //Not everyone can call anything
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Specialization string             `json:"specialization"`
	UseModel       primitive.ObjectID `json:"useModel" bson:"useModel"`
	CreatedAt      tools.Milliseconds `json:"createdAt"`
	UpdatedAt      tools.Milliseconds `json:"updatedAt"`
	CreatedBy      primitive.ObjectID `json:"createdBy"`
	UpdatedBy      primitive.ObjectID `json:"updatedBy"`
	Url            string
	XApiKey        primitive.ObjectID  `json:"xApiKey" bson:"xApiKey"`
	Args           []LLMAgentParameter `json:"args"`
	AgentProps     LLMAgentParams      `json:"agentProps" bson:"inline"`
}

type LLMAgentParams struct {
	Role               string               `json:"role"`
	SystemPrompt       string               `json:"systemPrompt" bson:"systemPrompt"`
	Tools              []string             `json:"tools"`
	RequiredAgents     []primitive.ObjectID `json:"requiredAgents"`
	AnswerFormat       string               `json:"answerFormat" bson:"answerFormat"`
	RequiredParameters []string             `json:"requiredParameters" bson:"requiredParameters"`
	OptionalParameters []string             `json:"optionalParameters" bson:"optionalParameters"`
}

type LLMAgentParameter struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Conditions bson.M `json:"conditions"`
}

type LLMAgentPerformance struct {
	Uid            primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	AgentUid       primitive.ObjectID `json:"agentUid" bson:"agentUid"`
	UserText       string             `json:"userText" bson:"userText"`
	UserEvaluation int8               `json:"userEvaluation" bson:"userEvaluation"`
	At             tools.Milliseconds `json:"at"`
}
