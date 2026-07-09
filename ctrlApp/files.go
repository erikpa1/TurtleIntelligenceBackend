package ctrlApp

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"turtle/core/files"
	"turtle/db"
	"turtle/modelsApp"
	"turtle/vfs"
)

const CT_FILES = "fileEntries"

// uploadSessions maps an in-progress upload's id to the virtual path it will
// be written to once all of its chunks have arrived.
var uploadSessions sync.Map

// cleanVirtualPath normalizes a slash-separated virtual path, stripping any
// leading/trailing slashes and collapsing "." to the root ("").
func cleanVirtualPath(relPath string) string {
	cleaned := strings.Trim(path.Clean("/"+relPath), "/")
	if cleaned == "." {
		return ""
	}
	return cleaned
}

// ensureVirtualFolder makes sure a folder row (and all of its ancestors)
// exists in the database, without touching the physical filesystem.
func ensureVirtualFolder(relPath string) {
	folderPath := cleanVirtualPath(relPath)
	if folderPath == "" {
		return
	}

	if db.EntityExists(CT_FILES, bson.M{"path": folderPath, "isDir": true}) {
		return
	}

	parent := cleanVirtualPath(path.Dir(folderPath))
	ensureVirtualFolder(parent)

	db.InsertEntity(CT_FILES, &modelsApp.FileEntry{
		Name:     path.Base(folderPath),
		Path:     folderPath,
		Parent:   parent,
		IsDir:    true,
		Modified: time.Now().Format(time.RFC3339),
	})
}

func QueryFiles(relPath string) []*modelsApp.FileEntry {
	entries := db.QueryEntities[modelsApp.FileEntry](CT_FILES, bson.M{"parent": cleanVirtualPath(relPath)})

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	folderPaths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir {
			folderPaths = append(folderPaths, entry.Path)
		}
	}

	counts := childCounts(folderPaths)
	for _, entry := range entries {
		if entry.IsDir {
			entry.Count = counts[entry.Path]
		} else {
			entry.Extension = strings.ToLower(files.ExtractExtensionNoDot(entry.Name))
		}
	}

	return entries
}

// childCounts aggregates, in a single query, how many direct children each
// of the given folder paths has.
func childCounts(folderPaths []string) map[string]int64 {
	result := map[string]int64{}
	if len(folderPaths) == 0 {
		return result
	}

	cursor, err := db.DB.Col(CT_FILES).Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"parent": bson.M{"$in": folderPaths}}}},
		{{Key: "$group", Value: bson.M{"_id": "$parent", "count": bson.M{"$sum": 1}}}},
	})
	if err != nil {
		return result
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var row struct {
			Parent string `bson:"_id"`
			Count  int64  `bson:"count"`
		}
		if err := cursor.Decode(&row); err != nil {
			continue
		}
		result[row.Parent] = row.Count
	}

	return result
}

// CreateFolder inserts a virtual folder row. It never creates a physical
// directory.
func CreateFolder(relPath string) error {
	folderPath := cleanVirtualPath(relPath)
	if folderPath == "" {
		return nil
	}

	if db.EntityExists(CT_FILES, bson.M{"path": folderPath}) {
		return nil
	}

	parent := cleanVirtualPath(path.Dir(folderPath))
	ensureVirtualFolder(parent)

	_, err := db.InsertEntity(CT_FILES, &modelsApp.FileEntry{
		Name:     path.Base(folderPath),
		Path:     folderPath,
		Parent:   parent,
		IsDir:    true,
		Modified: time.Now().Format(time.RFC3339),
	})
	return err
}

func ReadFile(relPath string) (string, error) {
	return db.SC.GetFileStringNew(cleanVirtualPath(relPath))
}

// ReadFileBytes is the binary counterpart of ReadFile, used to serve
// previews (images, PDFs, media, ...) without a text-encoding round trip.
func ReadFileBytes(relPath string) ([]byte, error) {
	return db.SC.GetFileBytesNew(cleanVirtualPath(relPath))
}

// WriteFile persists a text file's bytes through db.SC and upserts its
// metadata row, creating any missing virtual ancestor folders along the way.
func WriteFile(relPath string, content string) error {
	return WriteFileBytes(relPath, []byte(content))
}

// WriteFileBytes is the binary counterpart of WriteFile, used for uploads and
// anything else that isn't plain text.
func WriteFileBytes(relPath string, data []byte) error {
	filePath := cleanVirtualPath(relPath)

	parent := cleanVirtualPath(path.Dir(filePath))
	ensureVirtualFolder(parent)

	if _, err := db.SC.UploadFileNew(filePath, data); err != nil {
		return err
	}

	modified := time.Now().Format(time.RFC3339)
	size := int64(len(data))

	existing := db.QueryEntity[modelsApp.FileEntry](CT_FILES, bson.M{"path": filePath, "isDir": false})
	if existing != nil {
		return db.UpdateOneCustom(CT_FILES, bson.M{"_id": existing.Uid}, bson.M{"$set": bson.M{
			"size":     size,
			"modified": modified,
		}})
	}

	_, err := db.InsertEntity(CT_FILES, &modelsApp.FileEntry{
		Name:     path.Base(filePath),
		Path:     filePath,
		Parent:   parent,
		IsDir:    false,
		Size:     size,
		Modified: modified,
	})
	return err
}

// uploadScratchPath is where an in-progress upload's chunks are accumulated
// on local disk before being handed to db.SC as a single blob. This is
// implementation scratch space, not a virtual folder the user browses.
func uploadScratchPath(uploadId string) string {
	return filepath.Join(vfs.GetWorkingDirectory(), ".uploads", uploadId)
}

// StartUpload opens a new chunked-upload session for relPath and returns its
// id. Send the file's bytes to AppendUploadChunk in order, then call
// FinishUpload to persist it through db.SC.
func StartUpload(relPath string) (string, error) {
	filePath := cleanVirtualPath(relPath)
	if filePath == "" {
		return "", errors.New("invalid file path")
	}

	uploadId := primitive.NewObjectID().Hex()
	uploadSessions.Store(uploadId, filePath)

	scratch := uploadScratchPath(uploadId)
	if err := os.MkdirAll(filepath.Dir(scratch), 0755); err != nil {
		uploadSessions.Delete(uploadId)
		return "", err
	}
	if err := os.WriteFile(scratch, []byte{}, 0644); err != nil {
		uploadSessions.Delete(uploadId)
		return "", err
	}

	return uploadId, nil
}

// AppendUploadChunk appends the next chunk of bytes to an open upload session.
func AppendUploadChunk(uploadId string, chunk []byte) error {
	if _, ok := uploadSessions.Load(uploadId); !ok {
		return errors.New("unknown upload session")
	}

	f, err := os.OpenFile(uploadScratchPath(uploadId), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(chunk)
	return err
}

// FinishUpload assembles the session's accumulated chunks and persists them
// through WriteFileBytes, then discards the scratch file.
func FinishUpload(uploadId string) error {
	relPathAny, ok := uploadSessions.Load(uploadId)
	if !ok {
		return errors.New("unknown upload session")
	}
	relPath := relPathAny.(string)

	scratch := uploadScratchPath(uploadId)
	data, err := os.ReadFile(scratch)
	if err != nil {
		return err
	}

	if err := WriteFileBytes(relPath, data); err != nil {
		return err
	}

	os.Remove(scratch)
	uploadSessions.Delete(uploadId)
	return nil
}

// AbortUpload discards an in-progress upload session and its scratch file.
func AbortUpload(uploadId string) {
	if _, ok := uploadSessions.Load(uploadId); ok {
		os.Remove(uploadScratchPath(uploadId))
		uploadSessions.Delete(uploadId)
	}
}

// DeleteEntry removes a file or folder row. Folders cascade to their
// descendants; only the physical bytes of actual files are ever deleted.
func DeleteEntry(relPath string) error {
	entryPath := cleanVirtualPath(relPath)
	if entryPath == "" {
		return nil
	}

	entry := db.QueryEntity[modelsApp.FileEntry](CT_FILES, bson.M{"path": entryPath})
	if entry == nil {
		return nil
	}

	if entry.IsDir {
		for _, child := range db.QueryEntities[modelsApp.FileEntry](CT_FILES, bson.M{"parent": entryPath}) {
			if err := DeleteEntry(child.Path); err != nil {
				return err
			}
		}
	} else if err := db.SC.DeleteFileNew(entryPath); err != nil {
		return err
	}

	db.DeleteEntity(CT_FILES, bson.M{"_id": entry.Uid})
	return nil
}

// RenameEntry renames/moves a file or folder. Folder renames cascade to
// descendants purely as metadata updates; file renames also move the
// underlying bytes through db.SC.
func RenameEntry(oldRelPath, newRelPath string) error {
	oldPath := cleanVirtualPath(oldRelPath)
	newPath := cleanVirtualPath(newRelPath)

	entry := db.QueryEntity[modelsApp.FileEntry](CT_FILES, bson.M{"path": oldPath})
	if entry == nil {
		return nil
	}

	newParent := cleanVirtualPath(path.Dir(newPath))
	ensureVirtualFolder(newParent)

	if entry.IsDir {
		for _, child := range db.QueryEntities[modelsApp.FileEntry](CT_FILES, bson.M{"parent": oldPath}) {
			if err := RenameEntry(child.Path, path.Join(newPath, child.Name)); err != nil {
				return err
			}
		}
	} else if err := db.SC.MoveFileNew(oldPath, newPath); err != nil {
		return err
	}

	return db.UpdateOneCustom(CT_FILES, bson.M{"_id": entry.Uid}, bson.M{"$set": bson.M{
		"name":   path.Base(newPath),
		"path":   newPath,
		"parent": newParent,
	}})
}
