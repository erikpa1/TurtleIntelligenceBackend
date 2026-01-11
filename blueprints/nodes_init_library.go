package blueprints

import (
	"turtle/blueprints/handlers"
	"turtle/blueprints/handlers/filesystem"
	"turtle/blueprints/handlers/llm"
	. "turtle/blueprints/library"
)

func InitFilesystemNodes() {
	RegisterNodeHandler("writeToFile", filesystem.PlayWriteToFileNode)
	RegisterNodeHandler("loadFileString", filesystem.PlayLoadFileStringNode)
	RegisterNodeHandler("foreachFileInFolder", filesystem.PlayForeachFileInFolder)
}

func InitExcelNodes() {
	RegisterNodeHandler(WRITE_TO_EXCEL, handlers.PlayWriteExcel)
}

func InitTriggers() {
	//Todo toto prepisat do jednotlivych packagov
	RegisterNodeHandler(HTTP_TRIGGER, handlers.PlayHttpTriggerNode)
}

func InitAINodes() {
	RegisterNodeHandler(LLM_AGENT_NODE, llm.PlayLLMNode)
}

func InitSqliteNodes() {
	RegisterNodeHandler(WRITE_TO_SQLITE, handlers.PlayWriteSqlite)
}

func InitNodesLibrary() {
	InitFilesystemNodes()
	InitExcelNodes()
	InitTriggers()
	InitAINodes()
	InitSqliteNodes()
}
