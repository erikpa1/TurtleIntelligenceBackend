package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

// PointCloud is the metadata record for one uploaded point cloud. The raw
// uploaded file and every octree node's binary payload are persisted through
// db.SC, addressed by virtual paths - never stored in this document.
type PointCloud struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	SourcePath  string             `json:"sourcePath" bson:"sourcePath"`
	Extension   string             `json:"extension" bson:"extension"`
	Status      string             `json:"status" bson:"status"` // "processing" | "ready" | "error"
	Error       string             `json:"error" bson:"error"`
	HasColor    bool               `json:"hasColor" bson:"hasColor"`
	TotalPoints int64              `json:"totalPoints" bson:"totalPoints"`
	NodeCount   int                `json:"nodeCount" bson:"nodeCount"`
	MaxDepth    int                `json:"maxDepth" bson:"maxDepth"`
	BoundsMin   [3]float64         `json:"boundsMin" bson:"boundsMin"`
	BoundsMax   [3]float64         `json:"boundsMax" bson:"boundsMax"`
	Created     string             `json:"created" bson:"created"`
}
