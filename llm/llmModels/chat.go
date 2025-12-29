package llmModels

import (
	"fmt"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHistoryLight struct {
	Uid     primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	At      tools.Milliseconds `json:"at"`
	UserUid primitive.ObjectID `json:"userUid" bson:"userUid"`
	Org     primitive.ObjectID `json:"org"`
}

func ChatHistoryLightProjection() bson.M {
	return bson.M{
		"uid":  1,
		"name": 1,
		"at":   1,
	}

}

type ConversationSegment struct {
	At         tools.Milliseconds `json:"at"`
	Text       string             `json:"text"`
	IsUser     bool               `json:"isUser" bson:"isUser"`
	Duration   tools.Milliseconds `json:"duration"`
	SmartTexts []ContentBlock     `json:"smartTexts" bson:"smartTexts"`
}

type ChatHistory struct {
	ChatHistoryLight `bson:"inline"`
	Conversation     []ConversationSegment `json:"conversation"`
	Answered         bool                  `json:"answered"`
}

// ContentBlock represents a detected content block
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ChatRequestParams struct {
	SystemPrompt string
	UserPrompt   string
	Memory       string
}

func (self *ChatRequestParams) GetFinalCommand() string {

	resultBuffer := ""

	if self.Memory != "" {
		resultBuffer += fmt.Sprintf(`[PREVIOUS_MESSAGE]%s[/PREVIOUS_MESSAGE]`, self.Memory)
	}

	if self.SystemPrompt != "" {
		resultBuffer += fmt.Sprintf("[INST][SYSTEM] %s[/INST]", self.SystemPrompt)
	}

	if self.UserPrompt != "" {
		resultBuffer += fmt.Sprintf(`User query: %s`, self.UserPrompt)
	}

	return resultBuffer
}

func MCPSystemPrefix(text string) string {
	return fmt.Sprintf("You are MCP agent in the field of %s", text)
}
