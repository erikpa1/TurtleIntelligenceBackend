package agents

import (
	"fmt"

	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/node"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CT_AGENT_NODES = "llm_agent_nodes"
const CT_AGENT_EDGES = "llm_agent_edges"

func InsertNodes(user *models.User, nodes []*LLMAgentNode) {
	for _, n := range nodes {
		n.Org = user.Org
	}

	db.InsertMany(CT_AGENT_NODES, tools.ToIArray(nodes))

	lg.LogI(fmt.Sprintf("Inserted %d nodes", len(nodes)))

}

func COUNode(user *models.User, node *LLMAgentNode) {
	node.Org = user.Org

	if node.Uid.IsZero() {
		db.InsertEntity(CT_AGENT_NODES, node)
	} else {
		if user.Org == node.Org {
			db.SetById(CT_AGENT_NODES, node.Uid, node)
		} else {
			lg.LogE("Failed to update node for org %s", node.Org)
		}
	}
}

func GetAgentNode(orgUid, uid primitive.ObjectID) *LLMAgentNode {
	return db.GetByIdAndOrg[LLMAgentNode](CT_AGENT_NODES, uid, orgUid)
}

func GetRelationOfNode(orgUid primitive.ObjectID, query bson.M) *node.NodeRelation {
	query["org"] = orgUid
	return db.QueryEntity[node.NodeRelation](CT_AGENT_NODES, query)
}

func QueryNodes(user *models.User, query bson.M) []*LLMAgentNode {
	return db.QueryEntities[LLMAgentNode](CT_AGENT_NODES, user.FillOrgQuery(query))
}

func DeleteNodesOfAgent(user *models.User, agentUid primitive.ObjectID) {
	db.DeleteEntities(CT_AGENT_NODES, user.FillOrgQuery(bson.M{"parent": agentUid}))
	db.DeleteEntities(CT_AGENT_EDGES, user.FillOrgQuery(bson.M{"parent": agentUid}))
}

func DeleteAgentNode(nodeUid primitive.ObjectID) {
	//TODO musia sa zmazat secky connectiony

	db.DeleteEntities(CT_AGENT_EDGES, bson.M{"$or": bson.A{
		bson.M{"source": nodeUid},
		bson.M{"target": nodeUid},
	}})

	db.Delete(CT_AGENT_NODES, nodeUid)
}

func GetTargetOfNode(org primitive.ObjectID, uid primitive.ObjectID, connName string) *LLMAgentNode {

	edge := db.QueryEntity[NodeEdge](CT_AGENT_EDGES, bson.M{
		"source":       uid,
		"sourceHandle": connName,
	})

	if edge != nil {
		return GetAgentNode(edge.Org, edge.Uid)

	}
	return nil
}

func GetTargetsOfNode(context *NodePlayContext, uid primitive.ObjectID, connName string) []*LLMAgentNode {

	opts := options.FindOptions{
		Sort: bson.D{{Key: "priority", Value: 1}},
	}

	query := bson.M{
		"source": uid,
	}

	if connName != "" {
		query["sourceHandle"] = connName
	}

	edges := db.QueryEntities[NodeEdge](CT_AGENT_EDGES, query, &opts)

	nodes := make([]*LLMAgentNode, len(edges))

	for i, edge := range edges {
		_, isCycle := context.AlreadyPlayedNodes[edge.Target]
		if isCycle {
			lg.LogE("Cycle detected")
			break
		}
		nodes[i] = GetAgentNode(edge.Org, edge.Target)
	}

	return nodes
}

func PlayAgentNode(context *NodePlayContext, agentUid primitive.ObjectID) {

	entryNode := GetAgentNode(context.User.Org, agentUid)

	if entryNode != nil {
		DispatchPlayNode(context, entryNode)
	} else {
		lg.LogE("No node entry")
	}

}

func DispatchPlayNode(context *NodePlayContext, node *LLMAgentNode) {

	if node.PhaseType == AGENT_PHASE_TRIGGER {
		PlayHttpTriggerNode(context, node)
	} else {
		if node.Type == WRITE_TO_FILE {
			lg.LogI("Going to: ", node.Type)
			PlayWriteToFileNode(context, node)
		} else if node.Type == LLM_AGENT_NODE {
			lg.LogI("Going to: ", node.Type)
			PlayLLMNode(context, node)
		}
	}

	nextNodes := GetTargetsOfNode(context, node.Uid, "")

	if len(nextNodes) > 0 {

		for i, nextNode := range nextNodes {
			lg.LogI(fmt.Sprintf("[%d]-%s", i, nextNode.Name))
		}

		lg.LogOk("-----")

		for i, nextNode := range nextNodes {
			lg.LogI(fmt.Sprintf("[%d]-%s", i, nextNode.Name))

			if nextNode == nil {
				lg.LogE("No next node")
			} else {
				lg.LogE("Here")
				DispatchPlayNode(context, nextNode)
			}
		}
	} else {
		lg.LogE("No next nodes")
	}

}

func GetTypeDataOfNode[T any](parentUid primitive.ObjectID, conn string) *T {

	edgeQuery := bson.M{
		"sourceHandle": conn,
		"source":       parentUid,
	}

	edge := db.QueryEntity[NodeEdge](CT_AGENT_EDGES, edgeQuery)

	if edge != nil {
		node := db.QueryEntity[LLMAgentNode](CT_AGENT_NODES, bson.M{
			"_id": edge.Target,
		})

		if node != nil {
			return tools.RecastBson[T](node.TypeData)
		}
	}

	return nil
}
