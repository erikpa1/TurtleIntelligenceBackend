package documents

import (
	"github.com/gin-gonic/gin"
	"io"
	"turtle/auth"
	"turtle/lg"
	"turtle/tools"
)

func _ListDocuments(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListDocument(user))
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

func InitDocumentsApi(r *gin.Engine) {
	//Z nejakeho dovodu vo vite nefunguje /api/documents
	r.GET("/api/docs", auth.LoginOrApp, _ListDocuments)
	r.POST("/api/docs/upload", auth.LoginOrApp, _PostPdfDocument)
	r.DELETE("/api/docs/delete", auth.LoginOrApp, _PostPdfDocument)

}
