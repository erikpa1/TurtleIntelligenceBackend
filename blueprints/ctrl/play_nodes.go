package ctrl

import (
	"fmt"
	"turtle/blueprints/library"
	"turtle/blueprints/models"

	"turtle/lg"
)

func DispatchPlayNode(context *models.NodePlayContext, node *models.LLMAgentNode) {

	nodePlayFunc, nodePlayFuncExists := library.NODES_LIBRARY[node.Type]

	if nodePlayFuncExists {
		step := context.Pipeline.StartFromNode(node)
		nodePlayFunc(context, node)
		step.End()
	} else {
		lg.LogW("Unable to find node", node.Type)
	}

	nextNodes := GetTargetsOfNode(context, node.Uid, "")

	if len(nextNodes) > 0 {

		for i, nextNode := range nextNodes {
			lg.LogI(fmt.Sprintf("[%d]-%s", i, nextNode.Name))
		}

		lg.LogOk("-----")

		for _, nextNode := range nextNodes {

			if nextNode == nil {
				//lg.LogE("No next node")
			} else {
				DispatchPlayNode(context, nextNode)
			}
		}
	} else {
		//lg.LogE("No next nodes")
	}

}
