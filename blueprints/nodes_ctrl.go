package blueprints

import (
	"fmt"
	"turtle/blueprints/cts"
	"turtle/blueprints/library"
	"turtle/blueprints/models"
	"turtle/blueprints/utils"
	"turtle/core/users"
	"turtle/db"
	"turtle/lg"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNodes(user *users.User, nodes []*models.LLMAgentNode) {
	for _, n := range nodes {
		n.Org = user.Org
	}

	db.InsertMany(cts.CT_AGENT_NODES, tools.ToIArray(nodes))

	lg.LogI(fmt.Sprintf("Inserted %d nodes", len(nodes)))

}

func COUNode(user *users.User, node *models.LLMAgentNode) {
	node.Org = user.Org

	if node.Uid.IsZero() {
		db.InsertEntity(cts.CT_AGENT_NODES, node)
	} else {
		if user.Org == node.Org {
			db.SetById(cts.CT_AGENT_NODES, node.Uid, node)
		} else {
			lg.LogE("Failed to update node for org %s", node.Org)
		}
	}
}

func QueryNodes(user *users.User, query bson.M) []*models.LLMAgentNode {
	return db.QueryEntities[models.LLMAgentNode](cts.CT_AGENT_NODES, user.FillOrgQuery(query))
}

func DeleteNodesOfAgent(user *users.User, agentUid primitive.ObjectID) {
	db.DeleteEntities(cts.CT_AGENT_NODES, user.FillOrgQuery(bson.M{"parent": agentUid}))
	db.DeleteEntities(cts.CT_AGENT_EDGES, user.FillOrgQuery(bson.M{"parent": agentUid}))
}

func DeleteAgentNode(nodeUid primitive.ObjectID) {
	//TODO musia sa zmazat secky connectiony

	db.DeleteEntities(cts.CT_AGENT_EDGES, bson.M{"$or": bson.A{
		bson.M{"source": nodeUid},
		bson.M{"target": nodeUid},
	}})

	db.Delete(cts.CT_AGENT_NODES, nodeUid)
}

func PlayAgentNode(context *models.NodePlayContext, agentUid primitive.ObjectID) {

	entryNode := utils.GetAgentNode(context.User.Org, agentUid)

	if entryNode != nil {
		DispatchPlayNode(context, entryNode)
	} else {
		lg.LogE("No node entry")
	}

}

func DispatchPlayNode(context *models.NodePlayContext, node *models.LLMAgentNode) {

	nodePlayFunc, nodePlayFuncExists := library.NODES_LIBRARY[node.Type]

	if nodePlayFuncExists {
		nodePlayFunc(context, node)
	} else {
		lg.LogW("Unable to find node", node.Type)
	}

	nextNodes := utils.GetTargetsOfNode(context, node.Uid, "")

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
