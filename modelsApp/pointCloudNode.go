package modelsApp

import "go.mongodb.org/mongo-driver/bson/primitive"

// PointCloudNode is one octree node belonging to a PointCloud. Path is a
// digit string ("", "0", "03", ...) where each digit is the child octant
// index (0-7) at that depth; the root node has an empty Path. DataPath is
// the virtual path to the node's binary point payload, persisted via db.SC.
type PointCloudNode struct {
	Uid         primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	CloudUid    primitive.ObjectID `json:"cloudUid" bson:"cloudUid"`
	Path        string             `json:"path" bson:"path"`
	Depth       int                `json:"depth" bson:"depth"`
	PointCount  int                `json:"pointCount" bson:"pointCount"`
	BoundsMin   [3]float64         `json:"boundsMin" bson:"boundsMin"`
	BoundsMax   [3]float64         `json:"boundsMax" bson:"boundsMax"`
	HasChildren bool               `json:"hasChildren" bson:"hasChildren"`
	DataPath    string             `json:"dataPath" bson:"dataPath"`
}
