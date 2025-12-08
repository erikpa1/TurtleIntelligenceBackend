package agents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/node"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_AGENT_NODES = "llm_agent_nodes"

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

func QueryNodes(user *models.User, query bson.M) []bson.M {
	return db.QueryEntitiesAsCopy[bson.M](CT_AGENT_NODES, user.FillOrgQuery(query))
}

func DeleteNodesOfAgent(agentUid primitive.ObjectID) {
	db.DeleteEntities(CT_AGENT_NODES, bson.M{"parent": agentUid})
}

func DeleteAgentNode(nodeUid primitive.ObjectID) {
	db.Delete(CT_AGENT_NODES, nodeUid)
}
