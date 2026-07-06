package demand

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_DEMAND = "manufacturing_demand"

func ListDemand(user *users.User) []DemandOrder {
	return db.QueryEntitiesAsCopy[DemandOrder](CT_DEMAND, user.FillOrgQuery(bson.M{}))
}

func GetDemand(user *users.User, uid primitive.ObjectID) *DemandOrder {
	return db.GetByIdAndOrg[DemandOrder](CT_DEMAND, uid, user.Org)
}

func COUDemand(user *users.User, d *DemandOrder) {
	d.Org = user.Org

	if d.Uid.IsZero() {
		db.InsertEntity(CT_DEMAND, d)
	} else {
		db.SetByOrgAndId(CT_DEMAND, d.Uid, d.Org, d)
	}
}

func DeleteDemand(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_DEMAND, user.FillOrgQuery(bson.M{"_id": uid}))
}
