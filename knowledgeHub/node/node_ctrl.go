package node

import (
	"turtle/core/users"
	"turtle/db"
	"turtle/knowledgeHub/cts"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteNodesOfDomain(user *users.User, domainUid primitive.ObjectID) {
	db.DeleteEntities(cts.CT_KNOWLEDGE_NODES, bson.M{
		"org":    user.Org,
		"domain": domainUid,
	})
}
