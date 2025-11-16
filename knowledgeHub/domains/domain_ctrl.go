package domains

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/knowledgeHub/cts"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func COUDomain(user *models.User, domain *Domain) {
	domain.Org = user.Org
	db.COU(cts.CT_DOMAINS, domain.Uid, domain)
}

func DeleteDomain(user *models.User, uid primitive.ObjectID) {
	db.DeleteEntity(cts.CT_DOMAINS, bson.M{
		"org": user.Org,
		"_id": uid,
	})
}

func ListDomains(user *models.User) []*Domain {
	return db.QueryEntities[Domain](cts.CT_DOMAINS, bson.M{})
}
