package simulation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/modelsApp"
)

type SimEntity struct {
	Uid       primitive.ObjectID
	Name      string
	Type      string
	Behaviour ISimBehaviour
	Wold      *SimWorld
}

func (self *SimEntity) FromEntity(def *modelsApp.Entity) {
	self.Uid = def.Uid
	self.Name = def.Name
	self.Type = def.Type

	if self.Type == "spawn" {
		self.Behaviour = NewSpawnBehaviour()
	} else if self.Type == "process" {
		self.Behaviour = NewProcessBehaviour()
	} else {
		self.Behaviour = NewUndefinedBehaviour()
	}

	self.Behaviour.SetWorld(self.Wold)
}
