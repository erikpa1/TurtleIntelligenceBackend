package handlers

import (
	"io"
	"turtle/blueprints/models"
	"turtle/databases/tsqlite"
	"turtle/formats/office"

	"turtle/lg"
	"turtle/vfs"
)

func PlayHttpTriggerNode(context *models.NodePlayContext, node *models.LLMAgentNode) {

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

func PlayWriteExcel(context *models.NodePlayContext, node *models.LLMAgentNode) {
	step := context.Pipeline.StartFromNode(node)
	lg.LogE("Here")
	office.WriteExcel(vfs.GetWorkingDirectory() + "/Book1.xlsx")
	step.End()
}

func PlayWriteSqlite(context *models.NodePlayContext, node *models.LLMAgentNode) {
	step := context.Pipeline.StartFromNode(node)
	tsqlite.WriteJsonToSqlite(vfs.GetWorkingDirectory()+"/books.db", "books", "")
	step.End()
}
