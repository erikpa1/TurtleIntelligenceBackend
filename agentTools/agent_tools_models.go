package agentTools

import (
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentToolResult struct {
	TextRaw  string `json:"textRaw" bson:"textRaw"`
	TextInfo string `json:"textInfo" bson:"textInfo"`
	IsOk     bool   `json:"isOk" bson:"isOk"`
	IsDebug  bool   `json:"-" bson:"-"`
}

type AgentTool struct {
	Uid         primitive.ObjectID                         `json:"uid" bson:"_id,omitempty"`
	Name        string                                     `json:"name" bson:"name"`
	Description string                                     `json:"description" bson:"description"`
	Icon        string                                     `json:"icon" bson:"icon"`
	Inputs      string                                     `json:"inputs" bson:"inputs"`
	Type        string                                     `json:"type" bson:"type"`
	Provider    string                                     `json:"provider" bson:"provider"`
	Category    string                                     `json:"category" bson:"category"`
	Fn          func(result *AgentToolResult, data bson.M) `json:"-" bson:"-"`
}

func (self *AgentTool) CallFn(result *AgentToolResult, data bson.M) {
	defer tools.Recover(fmt.Sprintf("Failed to CALL [%s][%s]", self.Name, self.Uid.Hex()))

	if self.Fn != nil {
		defer self.Fn(result, data)
	} else {
		lg.LogE("Unable to call", self.Name, "FN is not defined")
	}

}

type AgentToolUsage struct {
	Uid        primitive.ObjectID `json:"uid"`
	Name       string             `json:"name"`
	Parameters bson.M             `json:"parameters"`
	ToolResult *AgentToolResult   `json:"toolResult"`
}

type AgentToolCall struct {
	SelectedTool primitive.ObjectID `json:"selected_tool"`
	Parameters   bson.M             `json:"parameters"`
	Reasoning    string             `json:"reasoning"`
}
