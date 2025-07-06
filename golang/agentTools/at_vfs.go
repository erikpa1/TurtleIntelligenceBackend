package agentTools

func InitVfsTools() {

	t1 := AgentTool{}
	t1.Uid = uid("vfsList")
	t1.Name = "Local FS List"
	t1.Description = "Iterates files on local disc"
	t1.Inputs = "filePath:string"
	t1.Icon = "share.svg"

	AGENT_TOOLS[t1.Uid] = &t1

}
