package simulation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/modelsApp"
)

type SimEntity struct {
	Uid  primitive.ObjectID
	Name string
	Type string
	Wold *SimWorld
}

func (self *SimEntity) FromEntity(def *modelsApp.Entity) {
	self.Uid = def.Uid
	self.Name = def.Name
	self.Type = def.Type

}
