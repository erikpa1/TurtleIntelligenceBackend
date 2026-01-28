package blueprints

import (
	. "turtle/blueprints/library"
	"turtle/blueprints/nodeHandlers"
	"turtle/blueprints/nodeHandlers/chat"
	"turtle/blueprints/nodeHandlers/filesystem"
	"turtle/blueprints/nodeHandlers/llm"
)

func InitFilesystemNodes() {
	RegisterNodeFunctionHandler("writeToFile", filesystem.PlayWriteToFileNode)
	RegisterNodeFunctionHandler("loadFileString", filesystem.PlayLoadFileStringNode)
	RegisterNodeFunctionHandler("foreachFileInFolder", filesystem.PlayForeachFileInFolder)
}

func InitExcelNodes() {
	RegisterNodeFunctionHandler(WRITE_TO_EXCEL, nodeHandlers.PlayWriteExcel)
}

func InitTriggers() {
	//Todo toto prepisat do jednotlivych packagov
	RegisterNodeFunctionHandler(HTTP_TRIGGER, nodeHandlers.PlayHttpTriggerNode)
	RegisterNodeFunctionHandler(CHAT_TRIGGER, chat.PlayChatTrigger)

}

func InitTests() {
	RegisterNodeFunctionHandler("testNode", chat.PlayChatTrigger)
}

func InitAINodes() {
	RegisterNodeFunctionHandler(LLM_AGENT_NODE, llm.PlayLLMNode)
}

func InitSqliteNodes() {
	RegisterNodeFunctionHandler(WRITE_TO_SQLITE, nodeHandlers.PlayWriteSqlite)
}

func InitNodesLibrary() {
	InitFilesystemNodes()
	InitExcelNodes()
	InitTriggers()
	InitAINodes()
	InitSqliteNodes()
}
