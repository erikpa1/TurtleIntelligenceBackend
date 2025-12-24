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
}

func (self *ChatRequestParams) GetFinalCommand() string {

	tmp := fmt.Sprintf(`
[INST]
[SYSTEM] %s
[/INST]
User query: %s
`, self.SystemPrompt, self.UserPrompt)

	return tmp
}
