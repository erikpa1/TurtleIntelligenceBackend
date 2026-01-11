package blueprints

import (
	"fmt"
	"turtle/blueprints/ctrl"
	"turtle/blueprints/cts"
	"turtle/blueprints/models"
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

	entryNode := ctrl.GetAgentNode(context.User.Org, agentUid)

	if entryNode != nil {
		ctrl.DispatchPlayNode(context, entryNode)
	} else {
		lg.LogE("No node entry")
	}

}
