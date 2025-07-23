package simulation

import (
	"fmt"
	"github.com/erikpa1/turtle/modelsApp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimActor struct {
	Name     string             `json:"name"`
	Id       int64              `json:"id"`
	Uid      primitive.ObjectID `json:"definition_uid"`
	Position [3]float32         `json:"position"`
	Color    string             `json:"color"`
	World    *SimWorld          `json:"-"`
}

func NewSimActor() *SimActor {
	tmp := &SimActor{}
	tmp.Position = [3]float32{0, 0, 0}
	return tmp
}

func (self *SimActor) FromActorDefinition(def *modelsApp.Actor) {
	self.Uid = def.Uid
	self.Name = fmt.Sprintf("%s_%d", def.Name, self.Id)
	self.Color = def.Color
}

func (self *SimActor) UpdatePosition(position [3]float32) {
	self.Position = position
	self.World.UpdateActorState(self.Id, "position", self.Position)
}
