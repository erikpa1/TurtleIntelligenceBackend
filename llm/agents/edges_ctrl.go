package agents

import (
	"fmt"
	"turtle/core/users"

	"turtle/db"
	"turtle/lg"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteEdges(user *users.User, query bson.M) {
	user.FillOrgQuery(query)
	db.DeleteEntities(CT_AGENT_EDGES, query)
}

func InsertEdges(user *users.User, edges []*NodeEdge) {
	for _, n := range edges {
		n.Org = user.Org
	}

	db.InsertMany(CT_AGENT_EDGES, tools.ToIArray(edges))
	lg.LogI(fmt.Sprintf("Inserted %d edges", len(edges)))
}

func QueryEdges(user *users.User, query bson.M) []*NodeEdge {
	return db.QueryEntities[NodeEdge](CT_AGENT_EDGES, user.FillOrgQuery(query))
}
