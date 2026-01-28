package ctrl

import (
	"fmt"
	"turtle/blueprints/library"
	"turtle/blueprints/models"

	"turtle/lgr"
)

func DispatchPlayNode(context *models.NodePlayContext, node *models.Node) {

	nodePlayFunc, nodePlayFuncExists := library.NODES_LIBRARY[node.Type]

	if nodePlayFuncExists {
		step := context.Pipeline.StartFromNode(node)
		nodePlayFunc(context, node)
		step.End()
	} else {
		lgr.Error("Unable to find node", node.Type)
	}

	nextNodes := GetTargetsOfNode(context, node.Uid, "")

	if len(nextNodes) > 0 {

		for i, nextNode := range nextNodes {
			lgr.Info(fmt.Sprintf("[%d]-%s", i, nextNode.Name))
		}

		lgr.Ok("-----")

		for _, nextNode := range nextNodes {

			if nextNode == nil {
				//lgr.Error("No next node")
			} else {
				DispatchPlayNode(context, nextNode)
			}
		}
	} else {
		//lgr.Error("No next nodes")
	}

}
