package blueprints

import (
	"errors"
	"fmt"
	"turtle/auth"
	"turtle/blueprints/ctrl"
	"turtle/blueprints/cts"
	"turtle/blueprints/models"
	"turtle/core/users"
	"turtle/db"
	"turtle/lgr"
	"turtle/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNodes(user *users.User, nodes []*models.Node) {
	for _, n := range nodes {
		n.Org = user.Org
	}

	db.InsertMany(cts.CT_AGENT_NODES, tools.ToIArray(nodes))

	lgr.Info(fmt.Sprintf("Inserted %d nodes", len(nodes)))

}

func COUNode(user *users.User, node *models.Node) {
	node.Org = user.Org

	if node.Uid.IsZero() {
		db.InsertEntity(cts.CT_AGENT_NODES, node)
	} else {
		if user.Org == node.Org {
			db.SetById(cts.CT_AGENT_NODES, node.Uid, node)
		} else {
			lgr.Error("Failed to update node for org %s", node.Org)
		}
	}
}

func QueryNodes(user *users.User, query bson.M) []*models.Node {
	return db.QueryEntities[models.Node](cts.CT_AGENT_NODES, user.FillOrgQuery(query))
}

func GetNode(orgUid primitive.ObjectID, uid primitive.ObjectID) *models.Node {
	return db.GetByIdAndOrg[models.Node](cts.CT_AGENT_NODES, uid, orgUid)
}

func GetBlueprintOfNode(user *users.User, userUid primitive.ObjectID) primitive.ObjectID {

	tmp := GetNode(user.Org, userUid)

	if tmp != nil {
		return tmp.Parent
	} else {
		return primitive.ObjectID{}
	}
}

func ListNodesOfBlueprintAsMap(user *users.User, parent primitive.ObjectID) map[primitive.ObjectID]*models.Node {
	nodes := map[primitive.ObjectID]*models.Node{}
	for _, node := range QueryNodes(user, bson.M{"parent": parent}) {
		nodes[node.Uid] = node
	}
	return nodes
}

func ListNodesOfBlueprint(user *users.User, parent primitive.ObjectID) []*models.Node {
	return QueryNodes(user, bson.M{"parent": parent})
}

func DeleteNodesOfBlueprint(user *users.User, agentUid primitive.ObjectID) {
	db.DeleteEntities(cts.CT_AGENT_NODES, user.FillOrgQuery(bson.M{"parent": agentUid}))
	db.DeleteEntities(cts.CT_AGENT_EDGES, user.FillOrgQuery(bson.M{"parent": agentUid}))
}

func DeleteNode(nodeUid primitive.ObjectID) {
	//TODO musia sa zmazat secky connectiony

	db.DeleteEntities(cts.CT_AGENT_EDGES, bson.M{"$or": bson.A{
		bson.M{"source": nodeUid},
		bson.M{"target": nodeUid},
	}})

	db.Delete(cts.CT_AGENT_NODES, nodeUid)
}

func PlayNode(context *models.NodePlayContext, agentUid primitive.ObjectID) {

	entryNode := ctrl.GetAgentNode(context.User.Org, agentUid)

	if entryNode != nil {
		ctrl.DispatchPlayNode(context, entryNode)
	} else {
		lgr.Error("No node entry")
	}

}

func PlayBlueprint(c *gin.Context, entryNode primitive.ObjectID) {
	user := auth.GetUserFromContext(c)

	parentBlueprint := GetBlueprintOfNode(user, entryNode)

	if parentBlueprint.IsZero() {
		tools.AutoErrorReturn(c, errors.New("not found blueprint"))
		return
	}

	nodes := ListNodesOfBlueprint(user, entryNode)

	if len(nodes) == 0 {
		tools.AutoErrorReturn(c, errors.New("not found nodes"))
		return
	}

	edges := ListEdgesOfBlueprint(user, parentBlueprint)

	if len(edges) == 0 {
		tools.AutoErrorReturn(c, errors.New("not found edges"))
		return
	}

	//playNodeContext := models.NodePlayContext{
	//	Gin:         c,
	//	User:        user,
	//	IsLocalHost: c.RemoteIP() == "::1",
	//	Nodes:       nodesMap,
	//}

}
