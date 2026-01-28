package chat

import (
	"turtle/blueprints/ctrl"
	"turtle/blueprints/models"
	"turtle/blueprints/types"
	"turtle/lgr"
	"turtle/tools"
)

type ChatTrigger struct {
	TriggerDescription string `json:"triggerDescription" bson:"triggerDescription"`
	ExamplePrompt      string `json:"examplePrompt" bson:"examplePrompt"`

	Context    *models.NodePlayContext `json:"-"`
	FlowOutput types.FlowOutput        `json:"-"`
}

func (self *ChatTrigger) Init(context *models.NodePlayContext) {
	self.Context = context
}

func (self *ChatTrigger) Call(inputName string) {
	self.FlowOutput.Call()
}

func PlayChatTrigger(context *models.NodePlayContext, node *models.Node) {

	data := tools.RecastBson[ChatTrigger](node.TypeData)

	if data != nil {

		nextNode := ctrl.GetTargetOfNode(context.User.Org, node.Uid, "end")

		if nextNode != nil {
			ctrl.DispatchPlayNode(context, nextNode)
		}
	} else {
		lgr.ErrorStack("Failed to cast node data to WriteToFileNode")
	}

}

func _ChatTriggerBuildCheck() {
	var _ models.INodeData = &ChatTrigger{}
}
