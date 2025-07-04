package agentTools

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/tools"
)

func uid(str string) primitive.ObjectID {
	tmp, err := tools.StringToObjectID(str)

	if err != nil {
		lg.LogE(err)
	}

	return tmp
}

func InitGoogleDiscTools() {

	t1 := AgentTool{}
	t1.Uid = uid("gDiscRead")
	t1.Name = "Google Disc File Read"
	t1.Description = "This tools enables to read from google disc file"
	t1.Icon = "gdisk.svg"
	t1.Inputs = "filePath:string"

	AGENT_TOOLS[t1.Uid] = &t1

	t2 := AgentTool{}
	t2.Uid = uid("gDiscWrite")
	t2.Name = "Google Disc File Write"
	t2.Description = "This tools enables to write to google drive"
	t2.Icon = "gdisk.svg"
	t1.Inputs = "filePath:string, fileData:string"

	AGENT_TOOLS[t2.Uid] = &t2

	t3 := AgentTool{}
	t3.Uid = uid("gDiscList")
	t3.Name = "Google Disc Files List"
	t3.Description = "This tools enables to write to google drive"
	t3.Icon = "gdisk.svg"

	AGENT_TOOLS[t3.Uid] = &t3

}
