package pointcloudApp

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"turtle/core/serverKit"
	"turtle/ctrlApp"
)

func _ListPointClouds(c *gin.Context) {
	serverKit.ReturnOkJson(c, ctrlApp.QueryPointClouds())
}

func _GetPointCloud(c *gin.Context) {
	uid, err := primitive.ObjectIDFromHex(c.Param("uid"))
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}

	cloud := ctrlApp.GetPointCloud(uid)
	if cloud == nil {
		serverKit.Return404(c, nil)
		return
	}
	serverKit.ReturnOkJson(c, cloud)
}

func _GetPointCloudTree(c *gin.Context) {
	uid, err := primitive.ObjectIDFromHex(c.Param("uid"))
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}
	serverKit.ReturnOkJson(c, ctrlApp.GetPointCloudTree(uid))
}

func _GetNodeData(c *gin.Context) {
	uid, err := primitive.ObjectIDFromHex(c.Param("uid"))
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}

	// The root node's Path is "", which can't survive as a URL segment (gin
	// collapses the resulting "//"), so callers send "root" instead.
	nodePath := c.Param("path")
	if nodePath == "root" {
		nodePath = ""
	}

	data, err := ctrlApp.GetNodeData(uid, nodePath)
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func _DeletePointCloud(c *gin.Context) {
	uid, err := primitive.ObjectIDFromHex(c.Param("uid"))
	if err != nil {
		serverKit.Return404(c, nil)
		return
	}
	if err := ctrlApp.DeletePointCloud(uid); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _StartUpload(c *gin.Context) {
	uploadId, err := ctrlApp.StartPointCloudUpload(c.PostForm("name"), c.PostForm("extension"))
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

	if err := ctrlApp.AppendPointCloudUploadChunk(c.PostForm("uploadId"), data); err != nil {
		serverKit.Return500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func _FinishUpload(c *gin.Context) {
	cloudUid, err := ctrlApp.FinishPointCloudUpload(c.PostForm("uploadId"))
	if err != nil {
		serverKit.Return500(c, err)
		return
	}
	serverKit.ReturnOkJson(c, gin.H{"cloudUid": cloudUid})
}

func _AbortUpload(c *gin.Context) {
	ctrlApp.AbortPointCloudUpload(c.PostForm("uploadId"))
	c.Status(http.StatusOK)
}

func InitPointCloudApi(r *gin.Engine) {
	r.GET("/api/pointclouds/list", _ListPointClouds)
	r.GET("/api/pointclouds/:uid", _GetPointCloud)
	r.GET("/api/pointclouds/:uid/tree", _GetPointCloudTree)
	r.GET("/api/pointclouds/:uid/node/:path/data", _GetNodeData)
	r.DELETE("/api/pointclouds/:uid", _DeletePointCloud)
	r.POST("/api/pointclouds/upload/start", _StartUpload)
	r.POST("/api/pointclouds/upload/chunk", _UploadChunk)
	r.POST("/api/pointclouds/upload/finish", _FinishUpload)
	r.POST("/api/pointclouds/upload/abort", _AbortUpload)
}
