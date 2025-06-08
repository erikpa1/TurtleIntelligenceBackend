package llmModels

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/tools"
)

type LLMAgent struct {
	Uid            primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org            primitive.ObjectID `json:"org" bson:"org"`
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
	CommandExample string              `json:"commandExample" bson:"commandExample"`
	UseReasoning   bool                `json:"useReasoning" bson:"useReasoning"`
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

type Mistral7bResponse struct {
	SelectedAgent primitive.ObjectID `json:"selected_agent" bson:"selected_agent"`
	Confidence    float32            `json:"confidence"`
	Parameters    bson.M             `json:"parameters"`
	Reasoning     string             `json:"reasoning"`
}

type AgentTestResponse struct {
	Uid       primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	AgentUid  primitive.ObjectID `json:"agentUid" bson:"agentUid"`
	Org       primitive.ObjectID `json:"org" bson:"org"`
	At        tools.Milliseconds `json:"at"`
	State     int8               `json:"state"`
	Result    Mistral7bResponse  `json:"result"`
	ResultRaw string             `json:"resultRaw"`
	Error     string             `json:"error"`
	Text      string             `json:"text"`
}
