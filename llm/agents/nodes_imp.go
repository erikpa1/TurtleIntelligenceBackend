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

	if node.Type == HTTP_TRIGGER {
		//TODO vybrat z body data

		bodyBytes, err := io.ReadAll(context.Gin.Request.Body)
		if err != nil {
			lg.LogStackTraceErr(err)
			return
		}

		// Convert bytes to string
		bodyString := string(bodyBytes)
		context.Data.SetString(bodyString)

		lg.LogE(bodyString)
	} else {
		lg.LogE("Undefined node")
	}

	step.End()

}

func PlayWriteToFileNode(context *NodePlayContext, node *LLMAgentNode) {

	step := context.Pipeline.NewStep()
	step.Start()

	data := tools.RecastBson[WriteToFileNode](node.TypeData)

	if data != nil {

		if data.UseWd {

			lg.LogE("Going to write", data.ParentFolder, data.GetFileName())
			lg.LogOk(context.Data.GetString())

			vfs.WriteFileStringToWD(data.ParentFolder, data.GetFileName(), context.Data.GetString())
		} else {
			vfs.WriteFileString(data.ParentFolder, data.GetFileName(), context.Data.GetString())
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

	data := GetTypeDataOfNode[OllamaNode](node.Uid, "llm")

	lg.LogEson(data)

	if data != nil {
		model := llmModels.LLM{}
		model.ModelVersion = data.ModelName

		lg.LogI("Going to chat with model:", data.ModelName)
		lg.LogI(context.Data.GetString())
		modelResponse := llmCtrl.ChatAgenticModelRaw(context.Gin, context.User, &model, context.Data.GetString())

		context.Data.SetString(modelResponse.ResultRaw)

	}

	step.End()
}
