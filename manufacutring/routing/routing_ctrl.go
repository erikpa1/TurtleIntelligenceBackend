package routing

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_ROUTINGS = "manufacturing_routings"

func ListRoutings(user *users.User) []Routing {
	return db.QueryEntitiesAsCopy[Routing](CT_ROUTINGS, user.FillOrgQuery(bson.M{}))
}

func GetRouting(user *users.User, uid primitive.ObjectID) *Routing {
	return db.GetByIdAndOrg[Routing](CT_ROUTINGS, uid, user.Org)
}

func COURouting(user *users.User, rt *Routing) {
	rt.Org = user.Org

	if rt.Operations == nil {
		rt.Operations = []Operation{}
	}
	// Normalise operation sequence to the submitted order.
	for i := range rt.Operations {
		rt.Operations[i].Sequence = i + 1
	}

	if rt.Uid.IsZero() {
		db.InsertEntity(CT_ROUTINGS, rt)
	} else {
		db.SetByOrgAndId(CT_ROUTINGS, rt.Uid, rt.Org, rt)
	}
}

func DeleteRouting(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_ROUTINGS, user.FillOrgQuery(bson.M{"_id": uid}))
}
