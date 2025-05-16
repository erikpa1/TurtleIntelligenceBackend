package tools

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AutoReturn(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
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
