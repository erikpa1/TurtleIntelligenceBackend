package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

// FileEntry is a virtual file or folder record. Folders exist only as
// database rows - no physical directory is ever created for them. Files
// carry the same metadata, with their bytes persisted through db.SC.
type FileEntry struct {
	Uid      primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Path     string             `json:"path" bson:"path"`
	Parent   string             `json:"parent" bson:"parent"`
	IsDir    bool               `json:"isDir" bson:"isDir"`
	Size     int64              `json:"size" bson:"size"`
	Modified string             `json:"modified" bson:"modified"`

	// Count is the number of direct children of a folder entry. It is
	// computed on read via aggregation, never stored.
	Count int64 `json:"count" bson:"-"`

	// Extension is derived from Name for file entries, never stored.
	Extension string `json:"extension" bson:"-"`
}
