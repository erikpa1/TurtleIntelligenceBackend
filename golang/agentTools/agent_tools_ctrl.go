package agentTools

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
)

var AGENT_TOOLS = make(map[primitive.ObjectID]*AgentTool)

func ListAgentTools() []*AgentTool {

	response := make([]*AgentTool, len(AGENT_TOOLS))

	index := 0
	for _, value := range AGENT_TOOLS {
		response[index] = value
		lg.LogOk(value.Name)
		lg.LogE(value.Uid.Hex())
		index++
	}

	return response

}

func InitCoreTools() {
	InitGoogleDiscTools()
	InitVfsTools()
}
