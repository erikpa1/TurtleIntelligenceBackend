package llm

import (
	"turtle/blueprints/ctrl"
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"
)

type LLMAgentData struct {
	SystemPrompt string `json:"systemPrompt" bson:"systemPrompt"`
}

func PlayLLMNode(context *models.NodePlayContext, node *models.Node) {

	llmData := ctrl.GetTypeDataOfNode[OllamaNode](node.Uid, "llm")
	myData := tools.RecastBson[LLMAgentData](node.TypeData)

	memory := ctrl.GetTargetOfNode(node.Org, node.Uid, "memory")

	if llmData != nil && myData != nil {

		model := llmModels.LLM{}
		model.ModelVersion = llmData.ModelName

		lg.LogI("Going to chat with model:", llmData.ModelName)
		lg.LogI(context.Data.GetString())

		chatRequest := llmModels.ChatRequestParams{}

		lg.LogEson(memory)

		if memory != nil {
			if memory.Type == "staticMemory" {
				staticMemData := tools.RecastBson[StaticMemory](memory.TypeData)
				if staticMemData != nil {
					chatRequest.Memory = staticMemData.MemoryText
				}
			}
		}

		chatRequest.UserPrompt = context.Data.GetString()
		chatRequest.SystemPrompt = myData.SystemPrompt

		lg.LogE(chatRequest.GetFinalCommand())

		modelResponse := llmCtrl.ChatModelWithSystem(context.Gin, context.User, &model, &chatRequest)

		context.Pipeline.ActiveStep.DataStr = modelResponse.ResultRaw
		context.Data.SetString(modelResponse.ResultRaw)
	}

}
