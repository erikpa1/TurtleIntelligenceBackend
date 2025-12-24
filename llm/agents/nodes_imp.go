package agents

import (
	"io"

	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"
	"turtle/vfs"
)

func PlayHttpTriggerNode(context *NodePlayContext, node *LLMAgentNode) {
	step := PipelineStep{
		Name: node.Name,
	}

	context.Pipeline.AddStep(&step)

	step.Start()

	bodyBytes, err := io.ReadAll(context.Gin.Request.Body)
	if err != nil {
		lg.LogStackTraceErr(err)
		return
	}

	// Convert bytes to string
	bodyString := string(bodyBytes)
	context.Data.SetString(bodyString)

	step.DataStr = bodyString
	step.End()

}

func PlayWriteToFileNode(context *NodePlayContext, node *LLMAgentNode) {

	step := context.Pipeline.NewStep()
	step.Start()

	data := tools.RecastBson[WriteToFileNode](node.TypeData)

	if data != nil {

		dataToWrite := context.Data.GetString()

		step.DataStr = dataToWrite

		if data.UseWd {
			//lg.LogE("Going to write", data.ParentFolder, data.GetFileName())
			//lg.LogOk(context.Data.GetString())
			vfs.WriteFileStringToWD(data.ParentFolder, data.GetFileName(), dataToWrite)
		} else {
			vfs.WriteFileString(data.ParentFolder, data.GetFileName(), dataToWrite)
		}

		if data.OpenFolder {
			if data.UseWd {
				vfs.OpenWDFolder(data.ParentFolder)
			} else {
				vfs.OpenFolder(data.ParentFolder)
			}

		}

	} else {
		lg.LogStackTraceErr("Failed to cast node data to WriteToFileNode")
	}

	lg.LogE(vfs.GetWorkingDirectory())

	step.End()
}

func PlayLLMNode(context *NodePlayContext, node *LLMAgentNode) {

	step := PipelineStep{
		Name: node.Name,
	}

	context.Pipeline.AddStep(&step)

	step.Start()

	llmData := GetTypeDataOfNode[OllamaNode](node.Uid, "llm")
	myData := tools.RecastBson[LLMAgentData](node.TypeData)
	lg.LogEson(llmData)

	if llmData != nil && myData != nil {

		model := llmModels.LLM{}
		model.ModelVersion = llmData.ModelName

		lg.LogI("Going to chat with model:", llmData.ModelName)
		lg.LogI(context.Data.GetString())
		modelResponse := llmCtrl.ChatModelWithSystem(context.Gin, context.User, &model, &llmModels.ChatRequestParams{
			SystemPrompt: myData.SystemPrompt,
			UserPrompt:   context.Data.GetString(),
		})

		step.DataStr = modelResponse.ResultRaw
		context.Data.SetString(modelResponse.ResultRaw)
	}

	step.End()
}
