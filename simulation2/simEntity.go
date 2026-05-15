package simulation2

import (
	"turtle/simulation/simMath"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"turtle/modelsApp"
	"turtle/tools"
)

type SimEntity struct {
	RuntimeId       int64
	Uid             primitive.ObjectID
	Name            string
	Type            string
	Position        simMath.Position
	TypeData        *tools.SafeJson
	World           *SimWorld
	Functions       map[string]any
	BehaviourEntity SimBehData
}

func (self *SimEntity) FromEntity(def *modelsApp.WorldEntity) {
	self.Uid = def.Uid
	self.Name = def.Name
	self.Type = def.Type
	self.Position = def.Position
	self.TypeData = tools.NewSafeJson()
	self.TypeData.Data = def.TypeData

}

func NewSimEntity() *SimEntity {
	return &SimEntity{
		Functions: make(map[string]any),
		TypeData:  tools.NewSafeJson(),
	}
}
