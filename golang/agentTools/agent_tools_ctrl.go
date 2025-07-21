package agentTools

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var AGENT_TOOLS = make(map[primitive.ObjectID]*AgentTool)

func ListAgentsTools() []*AgentTool {

	response := make([]*AgentTool, len(AGENT_TOOLS))

	index := 0
	for _, value := range AGENT_TOOLS {
		response[index] = value
		index++
	}

	return response

}

func GetAgentTool(toolsUid primitive.ObjectID) *AgentTool {
	return AGENT_TOOLS[toolsUid]
}

func GetToolsListForAgent(toolUids []primitive.ObjectID) []*AgentTool {

	response := make([]*AgentTool, 0)

	for _, value := range toolUids {

		tool, exists := AGENT_TOOLS[value]

		if exists {
			response = append(response, tool)
		}
	}

	return response

}

func InitCoreTools() {
	InitGoogleDiscTools()
	InitVfsTools()
}
