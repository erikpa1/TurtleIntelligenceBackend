package documents

import "C"
import (
	"fmt"
	"turtle/auth"
	"turtle/core/lgr"
	"turtle/core/serverKit"
	"turtle/tools"

	"github.com/gin-gonic/gin"
)

func _RefreshDocumentsCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	coll := GetDocumentCollection(user, uid)

	lgr.ErrorJson(uid)

	if coll == nil {
		lgr.Error(fmt.Sprintf("Document collection [%s] not found", uid))
		serverKit.Return404(c, "")
	} else {
		RefreshDocumentsCollection(c, user, coll)
	}

}

func _CreateDocumentsCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	data := tools.ObjFromJsonPtr[DocumentsCollection](c.PostForm("data"))
	CreateDocumentsCollection(c, user, data)
}

func _ListDocumentCollections(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	serverKit.ReturnOkJson(c, ListDocumentCollections(user))
}

func _ListDocumentsOfCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	serverKit.ReturnOkJson(c, ListDocumentsOfCollection(user, uid))
}

func _DeleteDocumentCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	DeleteDocumentsCollection(user, uid)
}

func _UnassignDocFromCollection(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	UnassignDocumentsCollection(user, uid)
}

func _ClearCollectionDocuments(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := serverKit.MongoObjectIdFromQuery(c)
	ClearDocumentsCollection(user, uid)
}

func InitDocumentCollectionsApi(r *gin.Engine) {
	r.GET("/api/docs-cols", auth.LoginOrApp, _ListDocumentCollections)
	r.GET("/api/docs-cols/docs", auth.LoginOrApp, _ListDocumentsOfCollection)
	r.POST("/api/docs-col", auth.LoginOrApp, _CreateDocumentsCollection)
	r.POST("/api/docs-col/refresh", auth.LoginOrApp, _RefreshDocumentsCollection)
	r.DELETE("/api/docs-col", auth.LoginOrApp, _DeleteDocumentCollection)
	r.DELETE("/api/docs-col/item", auth.LoginOrApp, _UnassignDocFromCollection)
	r.DELETE("/api/docs-col/clear", auth.LoginOrApp, _ClearCollectionDocuments)

}
