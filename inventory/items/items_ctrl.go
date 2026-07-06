package items

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_ITEMS = "inventory_items"

func ListItems(user *users.User) []Item {
	return db.QueryEntitiesAsCopy[Item](CT_ITEMS, user.FillOrgQuery(bson.M{}))
}

func GetItem(user *users.User, uid primitive.ObjectID) *Item {
	return db.GetByIdAndOrg[Item](CT_ITEMS, uid, user.Org)
}

// COUItem creates or updates an item, scoping it to the caller's organisation.
func COUItem(user *users.User, item *Item) {
	item.Org = user.Org

	if item.Uid.IsZero() {
		db.InsertEntity(CT_ITEMS, item)
	} else {
		db.SetByOrgAndId(CT_ITEMS, item.Uid, item.Org, item)
	}
}

func DeleteItem(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_ITEMS, user.FillOrgQuery(bson.M{"_id": uid}))
}
