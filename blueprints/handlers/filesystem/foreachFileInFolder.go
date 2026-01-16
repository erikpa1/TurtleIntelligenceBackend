package filesystem

import (
	"turtle/blueprints/ctrl"
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/tools"
	"turtle/vfs"
)

type ForeachFileInFolder struct {
	FolderPath string `json:"folderPath" bson:"folderPath"`
}

func (self *ForeachFileInFolder) GetFileName() string {
	if self.FolderPath == "" {
		return "./undefined"
	}
	return self.FolderPath
}

func PlayForeachFileInFolder(context *models.NodePlayContext, node *models.Node) {

	data := tools.RecastBson[ForeachFileInFolder](node.TypeData)

	if data != nil {

		files, err := vfs.ListFile(data.FolderPath)

		lg.LogE(files)

		if err == nil {
			for _, fileName := range files {
				context.Data.SetString(fileName)
				context.Pipeline.ActiveStep.DataStr = fileName

				nextNode := ctrl.GetTargetOfNode(context.User.Org, node.Uid, "loop")

				if nextNode != nil {
					ctrl.DispatchPlayNode(context, nextNode)
				}

			}
		} else {
			lg.LogE(files)
		}

		nextNode := ctrl.GetTargetOfNode(context.User.Org, node.Uid, "end")

		if nextNode != nil {
			ctrl.DispatchPlayNode(context, nextNode)
		}
	} else {
		lg.LogStackTraceErr("Failed to cast node data to WriteToFileNode")
	}

}
