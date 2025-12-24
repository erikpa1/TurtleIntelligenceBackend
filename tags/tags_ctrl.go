package tags

import (
	"turtle/core/users"
	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
)

const CT_TAGS = "tags"

func ListTags(user *users.User, query bson.M) []*Tag {
	return db.QueryEntities[Tag](CT_TAGS, user.FillOrgQuery(query))
}

func DeleteTag(user *users.User, tagUid string) {
	db.DeleteEntity(CT_TAGS, user.FillOrgQuery(bson.M{"_id": tagUid}))
}

func CreateTag(user *users.User, tag *Tag) {
	tag.Org = user.Org
	db.InsertEntity(CT_TAGS, tag)
}

func UpdateTag(user *users.User, tag *Tag) {
	db.UpdateOneCustom(CT_TAGS, user.FillOrgQuery(bson.M{
		"_id": tag.Uid,
	}), bson.M{"$set": tag})
}
