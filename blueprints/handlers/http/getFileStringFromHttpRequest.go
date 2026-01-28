package http

import (
	"io"
	"turtle/blueprints/models"
	"turtle/lgr"
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
			lgr.Error(err.Error())
			return
		}

		openedFile, err := file.Open()
		defer openedFile.Close()

		if err != nil {
			lgr.Error(err.Error())
			return
		}

		fileBytes, err := io.ReadAll(openedFile)

		if err != nil {
			lgr.Error(err.Error())
			return
		}

		// Convert bytes to string
		fileContent := string(fileBytes)
		data.BOutput.SetString(fileContent)

	} else {
		lgr.ErrorStack("Failed to cast node data to WriteToFileNode")
	}

}
