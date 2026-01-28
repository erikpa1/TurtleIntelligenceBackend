package types

import "turtle/blueprints/models"

//==============Output==============

type FlowOutput struct {
	Name           string `json:"name"`
	ConnectionNode models.INodeData
}

func (self *FlowOutput) Call() {

}

//==============Input==============

type FlowInput struct {
	Name string `json:"name"`
}
