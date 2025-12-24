package documents

import (
	"fmt"
	"turtle/core/users"
	"turtle/db"
	"turtle/llm/llmCtrl"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListDocumentCollections(user *users.User) []*DocumentsCollection {
	return db.QueryEntities[DocumentsCollection](CT_DOC_COLLECTION, user.FillOrgQuery(nil))
}

func ListDocumentsOfCollection(user *users.User, colUid primitive.ObjectID) []*Document {

	col := GetDocsCol(user, colUid)

	if col != nil {

		query := user.FillOrgQuery(bson.M{
			"_id": bson.M{
				"$in": col.Items,
			},
		})

		result := db.QueryEntities[Document](CT_DOC, query)
		return result
	} else {
		return make([]*Document, 0)
	}

}

func GetDocsCol(user *users.User, colUid primitive.ObjectID) *DocumentsCollection {
	query := user.FillOrgQuery(bson.M{
		"_id": colUid,
	})
	tmp := db.QueryEntity[DocumentsCollection](CT_DOC_COLLECTION, query)
	return tmp
}

func DeleteDocumentsCollection(user *users.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": uid,
	}))
}

func UnassignDocumentsCollection(user *users.User, uid primitive.ObjectID) {
	db.UpdateOneCustom(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": uid,
	}),
		bson.M{
			"$pull": bson.M{
				"items": uid, // the ObjectID you want to remove
			}},
	)
}

func ClearDocumentsCollection(user *users.User, uid primitive.ObjectID) {
	db.UpdateOneCustom(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": uid,
	}),
		bson.M{
			"$set": bson.M{
				"items": make([]primitive.ObjectID, 0),
			}},
	)
}

func GetDocumentCollection(user *users.User, id primitive.ObjectID) *DocumentsCollection {

	return db.QueryEntity[DocumentsCollection](CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": id,
	}))

}
func CreateDocumentsCollection(c *gin.Context, user *users.User, docColl *DocumentsCollection) {
	docColl.Org = user.Org
	docColl.Uid = primitive.NewObjectID()
	db.InsertEntity(CT_DOC_COLLECTION, docColl)
	RefreshDocumentsCollection(c, user, docColl)

}

func RefreshDocumentsCollection(c *gin.Context, user *users.User, docColl *DocumentsCollection) {

	ClearDocumentsCollection(user, docColl.Uid)

	for _, doc := range ListDocumentExtracts(user) {

		userQuery := fmt.Sprintf(`
	User command: %s
	Document text %s
			`,
			docColl.Filter,
			doc.Extraction)

		meetsRequirements := llmCtrl.AskTrueFalse(c, user, userQuery, 0.8)

		if meetsRequirements {
			db.UpdateOneCustom(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
				"_id": docColl.Uid,
			}),
				bson.M{
					"$push": bson.M{
						"items": doc.Uid,
					}},
			)
		}
	}

}
