package filesystem

import (
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/tools"
	"turtle/vfs"
)

func PlayLoadFileStringNode(context *models.NodePlayContext, node *models.LLMAgentNode) {
	step := context.Pipeline.StartFromNode(node)

	data := tools.RecastBson[LoadFileStringNode](node.TypeData)

	if data != nil {
		if data.UseWd {
			strData, err := vfs.GetFileStringFromWDNew(data.FilePath)
			context.Data.SetString(strData)
			step.SetError(err)
		} else {
			//lg.LogE("Going to write", data.ParentFolder, data.GetFileName())
			//lg.LogOk(context.Data.GetString())
			strData, err := vfs.GetFileString(data.FilePath)

			context.Data.SetString(strData)
			step.SetError(err)
		}

	} else {
		lg.LogStackTraceErr("Failed to cast node data to LoadFileStringNode")
	}

	step.End()
}

func PlayWriteToFileNode(context *models.NodePlayContext, node *models.LLMAgentNode) {

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
