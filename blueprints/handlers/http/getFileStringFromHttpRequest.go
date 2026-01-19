package http

import (
	"io"
	"turtle/blueprints/models"
	"turtle/lg"
	"turtle/tools"
)

type GetFileStringFromHttpRequest struct {
	FileFormName string `json:"folderPath" bson:"folderPath"`
	BOutput      models.ContextData
}

func PlayGetFileStringFromHttpRequest(context *models.NodePlayContext, node *models.Node) {

	data := tools.RecastBson[GetFileStringFromHttpRequest](node.TypeData)
	data.BOutput.SetString("")

	if data != nil {

		c := context.Gin

		file, err := c.FormFile("file")

		if err != nil {
			lg.LogE(err)
			return
		}

		openedFile, err := file.Open()
		defer openedFile.Close()

		if err != nil {
			lg.LogE(err)
			return
		}

		fileBytes, err := io.ReadAll(openedFile)

		if err != nil {
			lg.LogE(err)
			return
		}

		// Convert bytes to string
		fileContent := string(fileBytes)
		data.BOutput.SetString(fileContent)

	} else {
		lg.LogStackTraceErr("Failed to cast node data to WriteToFileNode")
	}

}
