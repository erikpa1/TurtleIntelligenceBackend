package domain

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledge/cts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func COUDomain(domain *Domain) {
	db.COU(cts.CT_DOMAINS, domain.Uid, domain)
}

func DeleteDomain(uid primitive.ObjectID) {
	db.Delete(cts.CT_DOMAINS, uid)
}

func ListDomains() []*Domain {
	return db.QueryEntities[Domain](cts.CT_DOMAINS, bson.M{})
}
