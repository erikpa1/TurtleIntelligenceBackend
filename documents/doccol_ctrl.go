package documents

import (
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/llm/llmCtrl"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListDocumentCollections(user *models.User) []*DocumentsCollection {
	return db.QueryEntities[DocumentsCollection](CT_DOC_COLLECTION, user.FillOrgQuery(nil))
}

func DeleteDocumentsCollection(user *models.User, uid primitive.ObjectID) {
	db.DeleteEntity(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": uid,
	}))
}

func GetDocumentCollection(user *models.User, id primitive.ObjectID) *DocumentsCollection {

	return db.QueryEntity[DocumentsCollection](CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": id,
	}))

}
func CreateDocumentsCollection(c *gin.Context, user *models.User, docColl *DocumentsCollection) {
	docColl.Org = user.Org
	docColl.Uid = primitive.NewObjectID()
	db.InsertEntity(CT_DOC_COLLECTION, docColl)
	RefreshDocumentsCollection(c, user, docColl)

}

func RefreshDocumentsCollection(c *gin.Context, user *models.User, docColl *DocumentsCollection) {

	db.UpdateOneCustom(CT_DOC_COLLECTION, user.FillOrgQuery(bson.M{
		"_id": docColl.Uid,
	}),
		bson.M{
			"$set": bson.M{
				"items": make([]primitive.ObjectID, 0),
			}},
	)

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
