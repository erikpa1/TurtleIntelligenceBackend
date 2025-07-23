package agentTools

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func InitVfsTools() {

	t1 := &AgentTool{}
	t1.Uid = uid("vfsList")
	t1.Name = "Filesystem List"
	t1.Description = "Iterates files on local disc"
	t1.Inputs = "filePath:string"
	t1.Icon = "share.svg"

	AGENT_TOOLS[t1.Uid] = t1

	t1 = &AgentTool{}
	t1.Uid = uid("vfsWrite")
	t1.Name = "Filesystem Write"
	t1.Description = "Writes to local file system"
	t1.Inputs = "filePath:string, fileBody:string"
	t1.Icon = "share.svg"
	t1.Fn = _VfsWrite

	AGENT_TOOLS[t1.Uid] = t1

}

func _VfsWrite(data bson.M) {

	safe := tools.SafeJson{}
	safe.Data = data

	lg.LogOkson(data)

	db.SC.UploadFileString("llm", safe.GetString("filePath", "x.txt"), safe.GetString("fileBody", "--not_found--"))

}
