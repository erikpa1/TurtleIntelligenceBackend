package agentTools

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/tools"
)

type AgentTool struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
	Inputs      string             `json:"inputType" bson:"inputType"`
	Fn          func(data bson.M)  `json:"-" bson:"-"`
}

func (self *AgentTool) CallFn(data bson.M) {
	defer tools.Recover(fmt.Sprintf("Failed to CALL [%s][%s]", self.Name, self.Uid.Hex()))

	if self.Fn != nil {
		defer self.Fn(data)
	} else {
		lg.LogE("Unable to call", self.Name, "FN is not defined")
	}

}
