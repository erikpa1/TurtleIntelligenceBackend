package handlers

import (
	"io"
	"turtle/blueprints/models"
	"turtle/databases/tsqlite"
	"turtle/formats/office"

	"turtle/lg"
	"turtle/vfs"
)

func PlayHttpTriggerNode(context *models.NodePlayContext, node *models.Node) {

	bodyBytes, err := io.ReadAll(context.Gin.Request.Body)
	if err != nil {
		lg.LogStackTraceErr(err)
		return
	}

	// Convert bytes to string
	bodyString := string(bodyBytes)
	context.Data.SetString(bodyString)
	context.Pipeline.ActiveStep.DataStr = bodyString

}

func PlayWriteExcel(context *models.NodePlayContext, node *models.Node) {
	office.WriteExcel(vfs.GetWorkingDirectory() + "/Book1.xlsx")
}

func PlayWriteSqlite(context *models.NodePlayContext, node *models.Node) {
	tsqlite.WriteJsonToSqlite(vfs.GetWorkingDirectory()+"/books.db", "books", "")
}
