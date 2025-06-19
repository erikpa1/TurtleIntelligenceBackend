package documents

import (
	"github.com/gin-gonic/gin"
	"io"
	"turtle/auth"
	"turtle/tools"
)

func _ListDocuments(c *gin.Context) {
	user := auth.GetUserFromContext(c)
	tools.AutoReturn(c, ListDocument(user))
}

func _PostPdfDocument(c *gin.Context) {

	fileIsBig := false

	//TODO upravit taktiez toto
	if fileIsBig {
		//DO nothing
	} else {

		file, _, err := c.Request.FormFile("pdf")
		if err != nil {
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

		InsertDocument(nil, data)

	}

}

func InitApi(r *gin.Engine) {
	group := r.Group("/api/documents")

	group.GET("/", auth.LoginOrApp, _ListDocuments)
	group.POST("/upload", auth.LoginOrApp, _PostPdfDocument)
	group.DELETE("/delete", auth.LoginOrApp, _PostPdfDocument)

}
