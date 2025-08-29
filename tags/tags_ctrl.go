package tags

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
)

const CT_TAGS = "tags"

func ListTags(user *models.User, query bson.M) []*Tag {
	return db.QueryEntities[Tag](CT_TAGS, user.FillOrgQuery(query))
}

func DeleteTag(user *models.User, tagUid string) {
	db.DeleteEntity(CT_TAGS, user.FillOrgQuery(bson.M{"_id": tagUid}))
}

func CreateTag(user *models.User, tag *Tag) {
	tag.Org = user.Org
	db.InsertEntity(CT_TAGS, tag)
}

func UpdateTag(user *models.User, tag *Tag) {
	db.UpdateOneCustom(CT_TAGS, user.FillOrgQuery(bson.M{
		"_id": tag.Uid,
	}), bson.M{"$set": tag})
}
