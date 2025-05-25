package simulation

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/modelsApp"
)

type SimActor struct {
	Name          string
	Id            int64
	DefinitionUid primitive.ObjectID
	Position      [3]float32
}

func NewSimActor() *SimActor {
	tmp := &SimActor{}
	tmp.Position = [3]float32{0, 0, 0}
	return tmp
}

func (self *SimActor) FromActorDefinition(def *modelsApp.Actor) {
	self.DefinitionUid = def.Uid
	self.Name = fmt.Sprintf("%s_%d", def.Name, self.Id)
}
