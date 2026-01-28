package flows

import (
	"turtle/blueprints/models"
	"turtle/blueprints/types"
)

const FLOW_TEST_NODE = "flowTestNode"

type FlowTestNode struct {
	FlowOutput types.FlowOutput `json:"-"`
	FlowInput  types.FlowInput  `json:"-"`
}

func (self *FlowTestNode) Init(context *models.NodePlayContext) {
	self.FlowOutput.Call()
}

func (self *FlowTestNode) Call(inputName string) {

}

func _TestNodeBuildCheck() {
	var _ models.INodeData = &FlowTestNode{}
}
