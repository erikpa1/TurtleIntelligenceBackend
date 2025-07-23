package agentTools

func InitMicrosoftTools() {

	t1 := &AgentTool{}
	t1.Uid = uid("microsoftMail")
	t1.Name = "Send mail"
	t1.Description = "Sends mail using Microsoft"
	t1.Inputs = "recepients:string[], header:string, body:string"
	t1.Icon = "shape_line.svg"
	t1.Fn = _MathMultiply

	AGENT_TOOLS[t1.Uid] = t1

}
