package node

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeRelation struct {
	Uid         primitive.ObjectID `json:"uid"`
	A           primitive.ObjectID `json:"a"`
	B           primitive.ObjectID `json:"b"`
	Type        string             `json:"relationType" bson:"relationType"`
	Description string             `json:"description"`
}
