package filesApp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"turtle/core/files"
	"turtle/core/serverKit"
	"turtle/ctrlApp"
)

var imageExtensions = map[string]bool{
	"png": true, "jpg": true, "jpeg": true, "gif": true,
	"webp": true, "bmp": true, "svg": true, "ico": true,
}

func _ListFiles(c *gin.Context) {
	serverKit.ReturnOkJson(c, ctrlApp.QueryFiles(c.Query("path")))
}

func _ReadFile(c *gin.Context) {
	relPath := c.Query("path")

	content, err := ctrlApp.ReadFile(relPath)
	if err != nil {
		serverKit.Return500(c, err)
		return
	}

	serverKit.ReturnOkJson(c, gin.H{
		"path":    relPath,
		"content": content,
	})
}

// _ReadFileRaw serves a file's raw bytes with a content type derived from its
// extension, so the frontend can preview images/PDFs/media inline instead of
// only reading them as text.
func _ReadFileRaw(c *gin.Context) {
	relPath := c.Query("path")

	data, err := ctrlApp.ReadFileBytes(relPath)
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}

	fileName := path.Base(relPath)
	ext := strings.ToLower(files.ExtractExtensionNoDot(fileName))

	switch {
	case ext == "pdf":
		serverKit.AutoPdfOrErrorNotFound(c, fileName, data, nil)
	case imageExtensions[ext]:
		serverKit.AutoImageOrErrorNotFound(c, fileName, data, nil)
	default:
		contentType := mime.TypeByExtension("." + ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileName))
		c.Data(http.StatusOK, contentType, data)
	}
}

func _CreateFolder(c *gin.Context) {
	if err := ctrlApp.CreateFolder(c.PostForm("path")); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _WriteFile(c *gin.Context) {
	if err := ctrlApp.WriteFile(c.PostForm("path"), c.PostForm("content")); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _DeleteEntry(c *gin.Context) {
	if err := ctrlApp.DeleteEntry(c.Query("path")); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _RenameEntry(c *gin.Context) {
	if err := ctrlApp.RenameEntry(c.PostForm("oldPath"), c.PostForm("newPath")); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _StartUpload(c *gin.Context) {
	uploadId, err := ctrlApp.StartUpload(c.PostForm("path"))
	if err != nil {
		serverKit.Return500(c, err)
		return
	}
	serverKit.ReturnOkJson(c, gin.H{"uploadId": uploadId})
}

func _UploadChunk(c *gin.Context) {
	fileHeader, err := c.FormFile("chunk")
	if err != nil {
		serverKit.Return500(c, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		serverKit.Return500(c, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		serverKit.Return500(c, err)
		return
	}

	if err := ctrlApp.AppendUploadChunk(c.PostForm("uploadId"), data); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _FinishUpload(c *gin.Context) {
	if err := ctrlApp.FinishUpload(c.PostForm("uploadId")); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _AbortUpload(c *gin.Context) {
	ctrlApp.AbortUpload(c.PostForm("uploadId"))
	c.Status(http.StatusOK)
}

func InitFilesApi(r *gin.Engine) {
	r.GET("/api/files/list", _ListFiles)
	r.GET("/api/files/read", _ReadFile)
	r.GET("/api/files/raw", _ReadFileRaw)
	r.POST("/api/files/folder", _CreateFolder)
	r.POST("/api/files/write", _WriteFile)
	r.DELETE("/api/files", _DeleteEntry)
	r.PUT("/api/files/rename", _RenameEntry)
	r.POST("/api/files/upload/start", _StartUpload)
	r.POST("/api/files/upload/chunk", _UploadChunk)
	r.POST("/api/files/upload/finish", _FinishUpload)
	r.POST("/api/files/upload/abort", _AbortUpload)
}
