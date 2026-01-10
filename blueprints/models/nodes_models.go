package models

import (
	"turtle/core/users"
	"turtle/lg"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LLMAgentNode struct {
	Uid         primitive.ObjectID            `json:"uid" bson:"_id,omitempty"`
	Name        string                        `json:"name"`
	Parent      primitive.ObjectID            `json:"parent"`
	Org         primitive.ObjectID            `json:"org"`
	PosX        float32                       `json:"posX" bson:"posX"`
	PosY        float32                       `json:"posY" bson:"posY"`
	Type        string                        `json:"type"`
	TypeData    bson.M                        `json:"typeData" bson:"typeData"`
	Connections map[string]primitive.ObjectID `json:"connections"` //Connections are deleted by editor and it has to modify nodes
}

type NodeEdge struct {
	Uid          primitive.ObjectID `json:"uid" bson:"_id,omitempty"`
	Parent       primitive.ObjectID `json:"parent"`
	Source       primitive.ObjectID `json:"source"`
	SourceHandle string             `json:"sourceHandle" bson:"sourceHandle"`
	Target       primitive.ObjectID `json:"target"`
	TargetHandle string             `json:"targetHandle" bson:"targetHandle"`
	Priority     int8               `json:"priority"`
	Org          primitive.ObjectID `json:"org"`
}

type _ContextDataType int8

type ContextDataTypeClass struct {
	Null   _ContextDataType
	String _ContextDataType
	Json   _ContextDataType
	Xml    _ContextDataType
}

var ContextDataType = ContextDataTypeClass{
	Null:   0,
	String: 1,
	Json:   2,
	Xml:    3,
}

type NodePlayContext struct {
	Gin                *gin.Context
	User               *users.User
	Data               ContextData                 `json:"data"`
	AlreadyPlayedNodes map[primitive.ObjectID]bool `json:"alreadyPlayed"`
	Pipeline           Pipeline                    `json:"pipeline"`
	IsLocalHost        bool
}

type ContextData struct {
	Data any              `json:"data"`
	Type _ContextDataType `json:"type"`
}

func (self *ContextData) GetString() string {
	if self.Data != nil {
		return self.Data.(string)
	} else {
		return ""

	}
}
func (self *ContextData) SetString(data string) {
	self.Type = ContextDataType.String
	self.Data = data

	lg.LogI("Setting stirng: ", data)

}
