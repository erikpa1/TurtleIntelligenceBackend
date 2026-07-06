package workcenters

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_WORKCENTERS = "manufacturing_workcenters"

func ListWorkCenters(user *users.User) []WorkCenter {
	return db.QueryEntitiesAsCopy[WorkCenter](CT_WORKCENTERS, user.FillOrgQuery(bson.M{}))
}

func GetWorkCenter(user *users.User, uid primitive.ObjectID) *WorkCenter {
	return db.GetByIdAndOrg[WorkCenter](CT_WORKCENTERS, uid, user.Org)
}

func COUWorkCenter(user *users.User, w *WorkCenter) {
	w.Org = user.Org

	if w.Uid.IsZero() {
		db.InsertEntity(CT_WORKCENTERS, w)
	} else {
		db.SetByOrgAndId(CT_WORKCENTERS, w.Uid, w.Org, w)
	}
}

func DeleteWorkCenter(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_WORKCENTERS, user.FillOrgQuery(bson.M{"_id": uid}))
}
