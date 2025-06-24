package tools

import (
	"bytes"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const CONTENT_TYPE_OCTET = "application/octet-stream"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_PDF = "application/pdf"

const CONTENT_TYPE_GZIP = "application/pdf"

const CONTENT_ENCODING_GZIP = "gzip"

func AutoReturn(c *gin.Context, data any) {
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.JSON(http.StatusOK, data)
	}

}

func AutoReturnOrError(c *gin.Context, err error, data any) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, data)
	}

}

func AutoNotFound(c *gin.Context, data any) {
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.JSON(http.StatusOK, data)
	}

}

func AutoErrorReturn(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, err.Error())
}

func AutoErrorsReturn(c *gin.Context, errors []error) {
	errorMessages := []string{}
	for _, err := range errors {
		errorMessages = append(errorMessages, err.Error())
	}
	c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
}

func ServerError(c *gin.Context, message any) {
	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(message)})
}

func BytesFromGinFormFile(c *gin.Context) ([]byte, *multipart.FileHeader, error) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return []byte{}, nil, err
	}

	// Open the file
	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return []byte{}, nil, err
	}
	defer fileData.Close()

	// Read the file content into memory
	data := make([]byte, file.Size)
	if _, err := fileData.Read(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return []byte{}, nil, err
	}

	return data, file, nil

}

func SSEHeaders(c *gin.Context) {
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Header().Set("Cache-Control", "no-cache")
}

func Return405WithError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func ReturnAccessDenied(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}

func WriteCookie(c *gin.Context, key string, value string) {
	c.SetCookie(key, value,
		int((time.Hour * 24 * 365 * 100).Seconds()),
		"/",   // path
		"",    // domain
		false, // secure
		false, // httpOnly
	)
}

func MongoObjectIdFromQuery(c *gin.Context) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(c.Query("uid"))
	if err != nil {
		return primitive.NilObjectID
	}
	return objectId
}

func MongoObjectIdFromQueryByKey(c *gin.Context, key string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(c.Query(key))
	if err != nil {
		return primitive.NilObjectID
	}
	return objectId
}

func AutoZipFileOrErrorNotFound(c *gin.Context, data []byte, err error, contentType string) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.Header("Content-Encoding", CONTENT_ENCODING_GZIP)

		c.Data(http.StatusOK, contentType, data)
	}
}

func AutoFileJsonOrErrorNotFound(c *gin.Context, data []byte, err error) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.Data(http.StatusOK, "application/json", data)
	}
}
func AutoFileOrErrorNotFound(c *gin.Context, data []byte, err error) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.Data(http.StatusOK, "application/octet-stream", data)
	}
}

func AutoPdfOrErrorNotFound(c *gin.Context, fileName string, data []byte, err error) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		c.Header("Content-Type", CONTENT_TYPE_PDF)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileName))

		c.Data(http.StatusOK, CONTENT_TYPE_PDF, data)
	}
}

func AutoImageOrErrorNotFound(c *gin.Context, fileName string, data []byte, err error) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {

		ext := ExtractExtensionNoDot(fileName)

		if ext == "svg" {
			c.Header("Content-Type", "image/svg+xml")
		} else {
			c.Header("Content-Type", "image/"+ext)
		}

		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileName))

		c.Data(http.StatusOK, CONTENT_TYPE_PDF, data)
	}
}

func AutoFileDownloadOrErrorNotFound(c *gin.Context, fileName string, data []byte, err error) {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
	} else {
		reader := bytes.NewReader(data)
		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Expires", "0")
		c.Header("Cache-Control", "must-revalidate")
		c.Header("Pragma", "public")
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
		c.DataFromReader(http.StatusOK, int64(len(data)), "application/octet-stream", reader, nil)
	}
}
