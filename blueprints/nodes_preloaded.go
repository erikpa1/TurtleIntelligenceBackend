package blueprints

import (
	"turtle/blueprints/models"

	"go.mongodb.org/mongo-driver/bson"
)

var PRELOADED_NODES = make(map[string]models.LLMAgentNode)

func LoadPreloadedNodes() {

	//TODO preload ndoes from database

	QueryNodes(nil, bson.M{"preload": true})
}
