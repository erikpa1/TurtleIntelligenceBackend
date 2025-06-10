package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"io"
)

func _PostPdfDocument(c *gin.Context) {

	fileIsBig := false

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

		// Create a bytes.Reader which implements io.ReaderAt
		reader := bytes.NewReader(data)

		// Now you can use it with your PDF function
		pdfReader, err := pdf.NewReader(reader, int64(len(data)))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create PDF reader"})
			return
		}

		var buf bytes.Buffer
		b, err := pdfReader.GetPlainText()
		if err != nil {
			panic(err)
		}
		buf.ReadFrom(b)
		content := buf.String()
		fmt.Println(content)
		// Work with your PDF...
	}

}

func InitDocumentsApi(r *gin.Engine) {

}
