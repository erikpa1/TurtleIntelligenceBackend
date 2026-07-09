package ctrlApp

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"turtle/db"
	"turtle/modelsApp"
	"turtle/pointcloud"
	"turtle/tools"
	"turtle/vfs"
)

const CT_POINTCLOUDS = "pointClouds"
const CT_POINTCLOUD_NODES = "pointCloudNodes"

type pointCloudUploadSession struct {
	name string
	ext  string
}

// pointCloudUploadSessions maps an in-progress point cloud upload's id to
// the name/extension it will be persisted under once all chunks arrive.
// Kept separate from the generic file-upload sessions in files.go since
// finishing one triggers point cloud processing, not a plain file write.
var pointCloudUploadSessions sync.Map

func pointCloudUploadScratchPath(uploadId string) string {
	return filepath.Join(vfs.GetWorkingDirectory(), ".pointcloud-uploads", uploadId)
}

func StartPointCloudUpload(name, extension string) (string, error) {
	if name == "" {
		return "", errors.New("name is required")
	}

	uploadId := primitive.NewObjectID().Hex()
	pointCloudUploadSessions.Store(uploadId, pointCloudUploadSession{name: name, ext: extension})

	scratch := pointCloudUploadScratchPath(uploadId)
	if err := os.MkdirAll(filepath.Dir(scratch), 0755); err != nil {
		pointCloudUploadSessions.Delete(uploadId)
		return "", err
	}
	if err := os.WriteFile(scratch, []byte{}, 0644); err != nil {
		pointCloudUploadSessions.Delete(uploadId)
		return "", err
	}

	return uploadId, nil
}

func AppendPointCloudUploadChunk(uploadId string, chunk []byte) error {
	if _, ok := pointCloudUploadSessions.Load(uploadId); !ok {
		return errors.New("unknown upload session")
	}

	f, err := os.OpenFile(pointCloudUploadScratchPath(uploadId), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(chunk)
	return err
}

// FinishPointCloudUpload assembles the uploaded bytes, persists the raw
// source file, creates a "processing" PointCloud record, and kicks off
// asynchronous parsing/octree-building. It returns as soon as the record
// exists - callers should poll GetPointCloud for status.
func FinishPointCloudUpload(uploadId string) (string, error) {
	sessionAny, ok := pointCloudUploadSessions.Load(uploadId)
	if !ok {
		return "", errors.New("unknown upload session")
	}
	session := sessionAny.(pointCloudUploadSession)

	scratch := pointCloudUploadScratchPath(uploadId)
	data, err := os.ReadFile(scratch)
	if err != nil {
		return "", err
	}

	cloud := &modelsApp.PointCloud{
		Name:      session.name,
		Extension: session.ext,
		Status:    "processing",
		Created:   time.Now().Format(time.RFC3339),
	}
	insertResult, err := db.InsertEntity(CT_POINTCLOUDS, cloud)
	if err != nil {
		return "", err
	}
	cloudUid := insertResult.InsertedID.(primitive.ObjectID)

	cloud.SourcePath = pointCloudRawPath(cloudUid, session.ext)
	if _, err := db.SC.UploadFileNew(cloud.SourcePath, data); err != nil {
		db.UpdateOneCustom(CT_POINTCLOUDS, bson.M{"_id": cloudUid}, bson.M{"$set": bson.M{
			"status": "error",
			"error":  err.Error(),
		}})
		os.Remove(scratch)
		pointCloudUploadSessions.Delete(uploadId)
		return "", err
	}
	db.UpdateOneCustom(CT_POINTCLOUDS, bson.M{"_id": cloudUid}, bson.M{"$set": bson.M{
		"sourcePath": cloud.SourcePath,
	}})

	go processPointCloud(cloudUid, data, session.ext)

	os.Remove(scratch)
	pointCloudUploadSessions.Delete(uploadId)

	return cloudUid.Hex(), nil
}

func AbortPointCloudUpload(uploadId string) {
	if _, ok := pointCloudUploadSessions.Load(uploadId); ok {
		os.Remove(pointCloudUploadScratchPath(uploadId))
		pointCloudUploadSessions.Delete(uploadId)
	}
}

func pointCloudRawPath(cloudUid primitive.ObjectID, ext string) string {
	return "/.pointclouds/raw/" + cloudUid.Hex() + "." + ext
}

func pointCloudNodeDataPath(cloudUid primitive.ObjectID, nodePath string) string {
	if nodePath == "" {
		nodePath = "root"
	}
	return "/.pointclouds/" + cloudUid.Hex() + "/nodes/" + nodePath + ".bin"
}

// processPointCloud parses the uploaded bytes, builds the LOD octree, and
// persists each node's binary payload and metadata row. Runs detached from
// the upload request.
func processPointCloud(cloudUid primitive.ObjectID, data []byte, extension string) {
	defer tools.Recover("Failed to process point cloud " + cloudUid.Hex())

	fail := func(err error) {
		db.UpdateOneCustom(CT_POINTCLOUDS, bson.M{"_id": cloudUid}, bson.M{"$set": bson.M{
			"status": "error",
			"error":  err.Error(),
		}})
	}

	points, hasColor, err := pointcloud.Parse(extension, data)
	if err != nil {
		fail(err)
		return
	}

	root := pointcloud.BuildOctree(points, pointcloud.DefaultMaxPointsPerNode, pointcloud.DefaultMaxDepth)
	nodes := pointcloud.Flatten(root)

	maxDepth := 0
	var totalPoints int64
	for _, n := range nodes {
		if n.Depth > maxDepth {
			maxDepth = n.Depth
		}
		totalPoints += int64(len(n.Points))

		dataPath := pointCloudNodeDataPath(cloudUid, n.Path)
		if _, err := db.SC.UploadFileNew(dataPath, encodeNodePayload(n.Points, hasColor)); err != nil {
			fail(err)
			return
		}

		_, err := db.InsertEntity(CT_POINTCLOUD_NODES, &modelsApp.PointCloudNode{
			CloudUid:    cloudUid,
			Path:        n.Path,
			Depth:       n.Depth,
			PointCount:  len(n.Points),
			BoundsMin:   n.Bounds.Min,
			BoundsMax:   n.Bounds.Max,
			HasChildren: n.HasChildren(),
			DataPath:    dataPath,
		})
		if err != nil {
			fail(err)
			return
		}
	}

	db.UpdateOneCustom(CT_POINTCLOUDS, bson.M{"_id": cloudUid}, bson.M{"$set": bson.M{
		"status":      "ready",
		"hasColor":    hasColor,
		"totalPoints": totalPoints,
		"nodeCount":   len(nodes),
		"maxDepth":    maxDepth,
		"boundsMin":   root.Bounds.Min,
		"boundsMax":   root.Bounds.Max,
	}})
}

// encodeNodePayload writes a lean binary blob: positions first (pointCount*3
// little-endian float32, interleaved x,y,z), then colors (pointCount*3
// uint8, interleaved r,g,b) only if hasColor. No header - callers already
// know the point count and hasColor from the node/cloud metadata.
func encodeNodePayload(points []pointcloud.Point, hasColor bool) []byte {
	posSize := len(points) * 3 * 4
	colSize := 0
	if hasColor {
		colSize = len(points) * 3
	}
	buf := make([]byte, posSize+colSize)

	for i, p := range points {
		o := i * 12
		binary.LittleEndian.PutUint32(buf[o:], math.Float32bits(p.X))
		binary.LittleEndian.PutUint32(buf[o+4:], math.Float32bits(p.Y))
		binary.LittleEndian.PutUint32(buf[o+8:], math.Float32bits(p.Z))
	}

	if hasColor {
		colorStart := posSize
		for i, p := range points {
			o := colorStart + i*3
			buf[o] = p.R
			buf[o+1] = p.G
			buf[o+2] = p.B
		}
	}

	return buf
}

func QueryPointClouds() []*modelsApp.PointCloud {
	return db.QueryEntities[modelsApp.PointCloud](CT_POINTCLOUDS, bson.M{})
}

func GetPointCloud(uid primitive.ObjectID) *modelsApp.PointCloud {
	return db.QueryEntity[modelsApp.PointCloud](CT_POINTCLOUDS, bson.M{"_id": uid})
}

func GetPointCloudTree(uid primitive.ObjectID) []*modelsApp.PointCloudNode {
	return db.QueryEntities[modelsApp.PointCloudNode](CT_POINTCLOUD_NODES, bson.M{"cloudUid": uid})
}

func GetNodeData(cloudUid primitive.ObjectID, nodePath string) ([]byte, error) {
	node := db.QueryEntity[modelsApp.PointCloudNode](CT_POINTCLOUD_NODES, bson.M{"cloudUid": cloudUid, "path": nodePath})
	if node == nil {
		return nil, errors.New("node not found")
	}
	return db.SC.GetFileBytesNew(node.DataPath)
}

func DeletePointCloud(uid primitive.ObjectID) error {
	cloud := GetPointCloud(uid)
	if cloud == nil {
		return nil
	}

	if cloud.SourcePath != "" {
		db.SC.DeleteFileNew(cloud.SourcePath)
	}

	for _, node := range GetPointCloudTree(uid) {
		if node.DataPath != "" {
			db.SC.DeleteFileNew(node.DataPath)
		}
	}

	db.DeleteEntities(CT_POINTCLOUD_NODES, bson.M{"cloudUid": uid})
	db.DeleteEntity(CT_POINTCLOUDS, bson.M{"_id": uid})
	return nil
}
