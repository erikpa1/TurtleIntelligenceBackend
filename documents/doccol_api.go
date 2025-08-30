package documents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
)

func _RefreshDocumentsCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	coll := GetDocumentCollection(user, uid)
	RefreshDocumentsCollection(c, user, coll)
}

func _CreateDocumentsCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[DocumentsCollection](c.PostForm("data"))
	CreateDocumentsCollection(c, user, data)
}

func _ListDocumentCollections(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListDocumentCollections(user))
}

func _DeleteDocumentCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteDocumentsCollection(user, uid)
}
func InitDocumentCollectionsApi(r *gin.Engine) {
	r.GET("/api/docs-cols", auth.LoginOrApp, _ListDocumentCollections)
	r.POST("/api/docs-col", auth.LoginOrApp, _CreateDocumentsCollection)
	r.POST("/api/docs-col/refresh", auth.LoginOrApp, _RefreshDocumentsCollection)
	r.DELETE("/api/docs-col", auth.LoginOrApp, _DeleteDocumentCollection)

}
