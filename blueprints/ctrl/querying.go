package ctrl

import (
	"turtle/blueprints/cts"
	"turtle/blueprints/models"
	"turtle/db"
	"turtle/knowledgeHub/node"
	"turtle/lg"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTypeDataOfNode[T any](parentUid primitive.ObjectID, conn string) *T {

	edgeQuery := bson.M{
		"sourceHandle": conn,
		"source":       parentUid,
	}

	edge := db.QueryEntity[models.NodeEdge](cts.CT_AGENT_EDGES, edgeQuery)

	if edge != nil {
		node := db.QueryEntity[models.Node](cts.CT_AGENT_NODES, bson.M{
			"_id": edge.Target,
		})

		if node != nil {
			return tools.RecastBson[T](node.TypeData)
		}
	}

	return nil
}

func GetRelationOfNode(orgUid primitive.ObjectID, query bson.M) *node.NodeRelation {
	query["org"] = orgUid
	return db.QueryEntity[node.NodeRelation](cts.CT_AGENT_NODES, query)
}

func GetAgentNode(orgUid, uid primitive.ObjectID) *models.Node {
	return db.GetByIdAndOrg[models.Node](cts.CT_AGENT_NODES, uid, orgUid)
}

func GetTargetOfNode(org primitive.ObjectID, uid primitive.ObjectID, connName string) *models.Node {

	edge := db.QueryEntity[models.NodeEdge](cts.CT_AGENT_EDGES, bson.M{
		"source":       uid,
		"sourceHandle": connName,
	})

	if edge != nil {
		return GetAgentNode(edge.Org, edge.Target)

	}
	return nil
}

func GetTargetsOfNode(context *models.NodePlayContext, uid primitive.ObjectID, connName string) []*models.Node {

	opts := options.FindOptions{
		Sort: bson.D{{Key: "priority", Value: 1}},
	}

	query := bson.M{
		"source": uid,
	}

	if connName != "" {
		query["sourceHandle"] = connName
	}

	edges := db.QueryEntities[models.NodeEdge](cts.CT_AGENT_EDGES, query, &opts)

	nodes := make([]*models.Node, len(edges))

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
