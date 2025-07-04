package agentTools

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var AGENT_TOOLS = make(map[primitive.ObjectID]*AgentTool)

func ListAgentTools() map[primitive.ObjectID]*AgentTool {

	return AGENT_TOOLS

}

func InitCoreTools() {
	InitGoogleDiscTools()
	InitVfsTools()
}
