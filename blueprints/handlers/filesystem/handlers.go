package filesystem

import (
	"turtle/blueprints/models"
	"turtle/lgr"
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
			//lgr.Error("Going to write", data.ParentFolder, data.GetFileName())
			//lgr.Ok(context.Data.GetString())
			strData, err := vfs.GetFileString(data.FilePath)

			context.Data.SetString(strData)
			context.Pipeline.ActiveStep.SetError(err)
		}

	} else {
		lgr.ErrorStack("Failed to cast node data to LoadFileStringNode")
	}

}
