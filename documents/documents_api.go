package documents

import (
	"fmt"
	"github.com/erikpa1/TurtleIntelligenceBackend/auth"
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"github.com/gin-gonic/gin"
	"io"
)

func _ListDocuments(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListDocuments(user))
}

func _GetDocument(c *gin.Context) {
	docUid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, GetDocument(user, docUid))

}

func _GetDocumentFile(c *gin.Context) {
	docUid := tools.MongoObjectIdFromQuery(c)
	user := auth.GetUserFromContext(c)

	doc := GetDocument(user, docUid)

	if doc != nil {

		file, err := db.SC.GetFileBytes("documents", doc.FileUidName())

		tools.AutoPdfOrErrorNotFound(c,
			fmt.Sprintf(doc.FileFullName(), doc.Name),
			file,
			err)

	} else {
		tools.AutoNotFound(c, nil)
	}

}

func _ListVSearchDocuments(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	query := c.Query("query")

	data, err := ListVSearchDocuments(c, user, query, 0.3)

	tools.AutoReturnOrError(c, err, data)
}

func _DeletePdfDocument(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	uid := tools.MongoObjectIdFromQuery(c)
	DeleteDocument(user, uid)
}

func _PostPdfDocument(c *gin.Context) {

	user := auth.GetUserFromContext(c)
	fileIsBig := false

	uploadParams := tools.ObjFromJson[InsertDocumentParams](c.PostForm("data"))

	//TODO upravit taktiez toto
	if fileIsBig {
		//DO nothing`
	} else {

		file, _, err := c.Request.FormFile("pdf")

		if err != nil {
			lg.LogE(err.Error())
			c.JSON(400, gin.H{"error": "Failed to get file"})
			return
		}
		defer file.Close()

		// Read the entire file into memory
		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to read file"})
			return
		}

		CreateAndUploadDocument(c, user, &uploadParams, data)
		lg.LogI("Uploaded document")

	}

}

func _UpdateDoc(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	document := tools.ObjFromJson[Document](c.PostForm("data"))
	UpdateDocument(user, &document)
}

func InitDocumentsApi(r *gin.Engine) {
	//Z nejakeho dovodu vo vite nefunguje /api/documents
	r.GET("/api/doc", auth.LoginOrApp, _GetDocument)
	r.GET("/api/doc/file", auth.LoginOrApp, _GetDocumentFile)
	r.GET("/api/docs", auth.LoginOrApp, _ListDocuments)
	r.GET("/api/doc/vsearch", auth.LoginOrApp, _ListVSearchDocuments)
	r.POST("/api/docs/upload", auth.LoginOrApp, _PostPdfDocument)
	r.PUT("/api/docs", auth.LoginOrApp, _UpdateDoc)
	r.DELETE("/api/docs", auth.LoginOrApp, _DeletePdfDocument)

}
