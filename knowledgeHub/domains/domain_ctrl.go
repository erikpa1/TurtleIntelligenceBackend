package domains

import (
	"turtle/core/users"
	"turtle/db"
	"turtle/knowledgeHub/cts"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func COUDomain(user *users.User, domain *Domain) {
	domain.Org = user.Org
	db.COU(cts.CT_DOMAINS, domain.Uid, domain)
}

func DeleteDomain(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(cts.CT_DOMAINS, bson.M{
		"org": user.Org,
		"_id": uid,
	})
}

func ListDomains(user *users.User) []*Domain {
	return db.QueryEntities[Domain](cts.CT_DOMAINS, bson.M{})
}
