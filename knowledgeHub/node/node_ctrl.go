package node

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/cts"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteNodesOfDomain(user *models.User, domainUid primitive.ObjectID) {
	db.DeleteEntities(cts.CT_KNOWLEDGE_NODES, bson.M{
		"org":    user.Org,
		"domain": domainUid,
	})
}
