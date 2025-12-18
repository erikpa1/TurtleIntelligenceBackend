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

func DeleteNodesOfAgent(agentUid primitive.ObjectID) {
	db.DeleteEntities(CT_AGENT_NODES, bson.M{"parent": agentUid})
	db.DeleteEntities(CT_AGENT_EDGES, bson.M{"parent": agentUid})
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

func PlayAgentNode(user *models.User, agentUid primitive.ObjectID) {

	entryNode := GetAgentNode(user.Org, agentUid)

	if entryNode != nil {
		DispatchPlayNode(entryNode)

	} else {
		lg.LogE("No node entry")
	}

}

func DispatchPlayNode(node *LLMAgentNode) {
	if node.PhaseType == AGENT_PHASE_TRIGGER {
		_PlayTriggerNode(node)
	} else if node.PhaseType == AGENT_PHASE_CONTROL {
		_PlayControlNode(node)
	} else if node.PhaseType == AGENT_PHASE_END {
		_PlayEndNode(node)
	}

}

func _PlayTriggerNode(node *LLMAgentNode) {
	if node.PhaseType == AGENT_PHASE_TRIGGER {
		targetNode := GetTargetOfNode(node.Org, node.Uid, "b")
		if targetNode != nil {
			DispatchPlayNode(targetNode)
		}
	}
}

func _PlayControlNode(node *LLMAgentNode) {
	if node.Type == "writeToFile" {

	} else if node.Type == "llmAgent" {

	} else if node.Type == "ollama" {

	}
}

func _PlayEndNode(node *LLMAgentNode) {

}
