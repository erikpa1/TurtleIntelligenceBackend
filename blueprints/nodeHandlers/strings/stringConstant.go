package strings

import (
	"turtle/blueprints/models"
	"turtle/blueprints/types"
)

type StringConstant struct {
	Value string `json:"value"`

	Context     *models.NodePlayContext `json:"-"`
	ValueOutput types.StringOutput      `json:"-"`
}

func (self *StringConstant) Init(context *models.NodePlayContext) {
	self.Context = context
	self.ValueOutput.SetData(self.Value)
}

func (self *StringConstant) Call(inputName string) {
	//pass

}

func _ChatTriggerBuildCheck() {
	var _ models.INodeData = &StringConstant{}
}
