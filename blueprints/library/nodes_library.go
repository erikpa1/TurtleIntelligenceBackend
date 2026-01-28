package library

import "turtle/blueprints/models"

type PlayNodeFunction func(context *models.NodePlayContext, node *models.Node)

var NODES_LIBRARY_HANDLERS = map[string]PlayNodeFunction{}

func RegisterNodeFunctionHandler(nodeType string, handler PlayNodeFunction) {
	NODES_LIBRARY_HANDLERS[nodeType] = handler
}

var NODES_LIBRARY = map[string]models.INodeData{}

func RegisterNodeData(nodeType string, handler models.INodeData) {
	NODES_LIBRARY[nodeType] = handler
}
