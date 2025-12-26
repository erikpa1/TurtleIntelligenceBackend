package agents

import (
	"io"
	"turtle/databases/tsqlite"
	"turtle/formats/office"

	"turtle/lg"
	"turtle/llm/llmCtrl"
	"turtle/llm/llmModels"
	"turtle/tools"
	"turtle/vfs"
)

func PlayHttpTriggerNode(context *NodePlayContext, node *LLMAgentNode) {

	step := context.Pipeline.StartFromNode(node)

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

	step := context.Pipeline.StartFromNode(node)

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

	step := context.Pipeline.StartFromNode(node)

	llmData := GetTypeDataOfNode[OllamaNode](node.Uid, "llm")
	myData := tools.RecastBson[LLMAgentData](node.TypeData)

	memory := GetTargetOfNode(node.Org, node.Uid, "memory")

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

		step.DataStr = modelResponse.ResultRaw
		context.Data.SetString(modelResponse.ResultRaw)
	}

	step.End()
}

func PlayWriteExcel(context *NodePlayContext, node *LLMAgentNode) {
	step := context.Pipeline.StartFromNode(node)
	lg.LogE("Here")
	office.WriteExcel(vfs.GetWorkingDirectory() + "/Book1.xlsx")
	step.End()
}

func PlayWriteSqlite(context *NodePlayContext, node *LLMAgentNode) {
	step := context.Pipeline.StartFromNode(node)
	tsqlite.WriteJsonToSqlite(vfs.GetWorkingDirectory()+"/books.db", "books", "")
	step.End()
}
