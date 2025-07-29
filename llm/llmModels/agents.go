package llmModels

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/agentTools"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	XApiKey        primitive.ObjectID   `json:"xApiKey" bson:"xApiKey"`
	Args           []LLMAgentParameter  `json:"args"`
	AgentProps     LLMAgentParams       `json:"agentProps" bson:"inline"`
	CommandExample string               `json:"commandExample" bson:"commandExample"`
	UseReasoning   bool                 `json:"useReasoning" bson:"useReasoning"`
	IsConfidential bool                 `json:"isConfidential" bson:"isConfidential"` //Confidential znamen ze nemoze logovat, obsahuje data
	Tools          []primitive.ObjectID `json:"tools"`
}

type LLMAgentTool struct {
	Uid   primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Org   primitive.ObjectID `json:"org" bson:"org"`
	Agent primitive.ObjectID `json:"agent" bson:"agent"`
	Tool  primitive.ObjectID `json:"tool" bson:"tool"`
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
	Uid              primitive.ObjectID           `json:"uid" bson:"_id,omitempty"`
	AgentName        string                       `json:"agentName" bson:"agentName"`
	AgentUid         primitive.ObjectID           `json:"agentUid" bson:"agentUid"`
	ResponseAgentUid primitive.ObjectID           `json:"responseAgentUid" bson:"responseAgentUid"`
	Org              primitive.ObjectID           `json:"org"`
	At               tools.Milliseconds           `json:"at"`
	State            int8                         `json:"state"`
	Result           Mistral7bResponse            `json:"result"`
	ResultRaw        string                       `json:"resultRaw"`
	Error            string                       `json:"error"`
	Text             string                       `json:"text"`
	AgentToolUsage   []*agentTools.AgentToolUsage `json:"agentToolUsage"`
}

func NewAgentTestResponse() *AgentTestResponse {
	return &AgentTestResponse{
		AgentToolUsage: make([]*agentTools.AgentToolUsage, 0),
	}
}
