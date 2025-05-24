package simulation

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/tools"
)

type SpawnBehaviour struct {
	World  *SimWorld
	Entity SimEntity

	NextSpawnTime tools.Seconds
}

func NewSpawnBehaviour() *SpawnBehaviour {
	return &SpawnBehaviour{}
}

func (self *SpawnBehaviour) SpawnActorOfType(actorUid primitive.ObjectID) {
	self.World.SpawnActorWithUid(actorUid)
}

// Implementacia IBehaviour
func (self *SpawnBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *SpawnBehaviour) Step() {
	lg.LogI("Spawn is doing step, next spawn time:", self.NextSpawnTime)

}

func (self *SpawnBehaviour) Init1() {

}

func (self *SpawnBehaviour) Init2() {

}

func (self *SpawnBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = *entity
}
