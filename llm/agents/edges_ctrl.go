package agents

import (
	"fmt"

	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertEdges(user *models.User, edges []*LLMAgentConnection) {
	for _, n := range edges {
		n.Org = user.Org
	}

	db.InsertMany(CT_AGENT_EDGES, tools.ToIArray(edges))
	lg.LogI(fmt.Sprintf("Inserted %d edges", len(edges)))
}

func QueryEdges(user *models.User, query bson.M) []*LLMAgentConnection {
	return db.QueryEntities[LLMAgentConnection](CT_AGENT_EDGES, user.FillOrgQuery(query))
}
