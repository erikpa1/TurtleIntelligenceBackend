package filesystem

import (
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/tools"
	"turtle/vfs"
)

func PlayLoadFileStringNode(context *models.NodePlayContext, node *models.Node) {

	data := tools.RecastBson[LoadFileStringNode](node.TypeData)

	if data != nil {
		if data.UseWd {
			strData, err := vfs.GetFileStringFromWDNew(data.FilePath)
			context.Data.SetString(strData)
			context.Pipeline.ActiveStep.SetError(err)
		} else {
			//lg.LogE("Going to write", data.ParentFolder, data.GetFileName())
			//lg.LogOk(context.Data.GetString())
			strData, err := vfs.GetFileString(data.FilePath)

			context.Data.SetString(strData)
			context.Pipeline.ActiveStep.SetError(err)
		}

	} else {
		lg.LogStackTraceErr("Failed to cast node data to LoadFileStringNode")
	}

}
