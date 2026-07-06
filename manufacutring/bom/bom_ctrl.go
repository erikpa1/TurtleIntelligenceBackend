package bom

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_BOMS = "manufacturing_boms"

func ListBoms(user *users.User) []Bom {
	return db.QueryEntitiesAsCopy[Bom](CT_BOMS, user.FillOrgQuery(bson.M{}))
}

func GetBom(user *users.User, uid primitive.ObjectID) *Bom {
	return db.GetByIdAndOrg[Bom](CT_BOMS, uid, user.Org)
}

// COUBom creates or updates a bill of materials, scoping it to the caller's
// organisation. Components with no explicit position keep their submitted order.
func COUBom(user *users.User, b *Bom) {
	b.Org = user.Org

	if b.Components == nil {
		b.Components = []BomComponent{}
	}

	for i := range b.Components {
		b.Components[i].Position = i
	}

	if b.Uid.IsZero() {
		db.InsertEntity(CT_BOMS, b)
	} else {
		db.SetByOrgAndId(CT_BOMS, b.Uid, b.Org, b)
	}
}

func DeleteBom(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_BOMS, user.FillOrgQuery(bson.M{"_id": uid}))
}
