package agentTools

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func uid(str string) primitive.ObjectID {
	tmp, err := tools.StringToObjectID(str)

	if err != nil {
		lg.LogE(err)
	}

	return tmp
}

func InitGoogleTools() {

	t1 := &AgentTool{}
	t1.Uid = uid("gDiscRead")
	t1.Name = "Google Disc File Read"
	t1.Description = "This tools enables to read from google disc file"
	t1.Icon = "google-docs.svg"
	t1.Inputs = "filePath:string"

	AGENT_TOOLS[t1.Uid] = t1

	t1 = &AgentTool{}
	t1.Uid = uid("gDiscWrite")
	t1.Name = "Google Disc File Write"
	t1.Description = "This tools enables to write to google drive"
	t1.Icon = "google-docs.svg"
	t1.Inputs = "filePath:string, fileData:string"

	AGENT_TOOLS[t1.Uid] = t1

	t1 = &AgentTool{}
	t1.Uid = uid("gDiscList")
	t1.Name = "Google Disc Files List"
	t1.Description = "This tools enables to write to google drive"
	t1.Icon = "google-docs.svg"

	AGENT_TOOLS[t1.Uid] = t1

}
