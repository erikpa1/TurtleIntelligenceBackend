package library

import "turtle/blueprints/models"

type PlayNodeFunction func(context *models.NodePlayContext, node *models.LLMAgentNode)

var NODES_LIBRARY = map[string]PlayNodeFunction{}

func RegisterNodeHandler(nodeType string, handler PlayNodeFunction) {
	NODES_LIBRARY[nodeType] = handler
}
