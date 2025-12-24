package agentTools

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"turtle/tools"
)

func InitMathTools() {

	t1 := &AgentTool{}
	t1.Uid = uid("mathMultiply")
	t1.Name = "Multiply"
	t1.Description = "Multiplies 2 numbers"
	t1.Inputs = "numberA:float64, numberB:float64"
	t1.Icon = "shape_line.svg"
	t1.Fn = _MathMultiply

	AGENT_TOOLS[t1.Uid] = t1

}

func _MathMultiply(result *AgentToolResult, data bson.M) {

	safe := tools.SafeJson{}
	safe.Data = data

	numberA := safe.GetDouble("numberA", 0)
	numberB := safe.GetDouble("numberB", 0)

	resultNumber := numberA * numberB

	result.TextRaw = fmt.Sprintf("%d", resultNumber)
	result.TextInfo = fmt.Sprintf("Multiply value: %s", resultNumber)

}
