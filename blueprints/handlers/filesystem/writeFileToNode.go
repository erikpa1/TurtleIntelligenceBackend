package filesystem

import (
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/tools"
	"turtle/vfs"
)

type WriteToFileNode struct {
	ParentFolder string `json:"parentFolder" bson:"parentFolder"`
	FileName     string `json:"fileName" bson:"fileName"`
	OpenFolder   bool   `json:"openFolder" bson:"openFolder"`
	UseWd        bool   `json:"useWd" bson:"useWd"`
}

func (self *WriteToFileNode) GetFileName() string {
	if self.FileName == "" {
		return "output.txt"
	}
	return self.FileName
}

func PlayWriteToFileNode(context *models.NodePlayContext, node *models.LLMAgentNode) {

	data := tools.RecastBson[WriteToFileNode](node.TypeData)

	if data != nil {

		dataToWrite := context.Data.GetString()

		context.Pipeline.ActiveStep.DataStr = dataToWrite

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
}
